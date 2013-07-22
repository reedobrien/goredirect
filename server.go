// Copyright 2013 Reed O'Brien. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/howeyc/fsnotify"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

var rules map[string]map[string]string
var size int = 19
var status int = http.StatusNotFound
var user string = "-"
var VERSION string = "1.0.0"

func main() {
	address := flag.String("address", "127.0.0.1", "The address to listen on")
	rulesPath := flag.String("rules", "", "Path to the JSON file of redirects")
	port := flag.String("port", "8080", "The port to listen on")
	watch := flag.Bool("watch", false, "Watch for JSON rules file changes")
	flag.Parse()
	if *rulesPath == "" {
		log.Fatalln("You must supply a mapping file")
	}

	fmt.Printf("Starting godir...listening on http://%s:%s\n", *address, *port)

	log.Println("Loading rules from:", *rulesPath)
	err := loadRules(*rulesPath)
	if err != nil {
		log.Fatalln(err)
	}

	if *watch {
		log.Println("Starting watcher for", *rulesPath)
		watcher, err := fsnotify.NewWatcher()
		watcher.Watch(*rulesPath)
		if err != nil {
			log.Println("Can't watch file", *rulesPath)
		}
		go func() {
			for {
				select {
				case ev := <-watcher.Event:
					if ev.IsModify() {
						log.Println(ev.Name, "updated, attempting reload...")
						err := loadRules(ev.Name)
						if err != nil {
							log.Println("Couln't reload rules: ", err)
						} else {
							log.Println("Reloaded rules from", ev.Name)
						}
					}
				case err := <-watcher.Error:
					log.Println("Error watching file:", *rulesPath, err)
				}
			}
		}()
	}

	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(*address+":"+*port, Log(http.DefaultServeMux)))
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Server", fmt.Sprintf("goredirect/%s", VERSION))
	target := rules[strings.Split(r.Host, ":")[0]][r.URL.Path]
	if target == "" {
		http.NotFound(w, r)
		return
	}
	status = http.StatusMovedPermanently
	//  len("<a href="">Moved Permanently</a>." == 33
	size = len(target) + 33
	http.Redirect(w, r, target, status)
	return
}

func Log(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t := time.Now()
		if r.URL.User != nil {
			user = r.URL.User.Username()
		} else {
			user = "-"
		}
		fmt.Printf("%s - %s [%s] \"%s %s %s\" %d %d \"%s\" \"%s\"\n",
			strings.Split(r.RemoteAddr, ":")[0], user, t.Format("02/Jan/2006:15:04:05 -0700"), r.Method, r.URL, r.Proto, status, size, r.Referer(), r.UserAgent())
		handler.ServeHTTP(w, r)
	})
}

func loadRules(path string) (err error) {
	f, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	err = json.Unmarshal(f, &rules)
	if err != nil {
		return err
	}
	return
}
