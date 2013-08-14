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
	"net/url"
	"strings"
	"time"
)

const (
	NotFoundLength = 19
	RedirectLength = 35
)

var remoteAddr string
var rules map[string]map[string]map[string]string
var size int
var status int = http.StatusNotFound
var target map[string]string
var user string = "-"
var VERSION string = "2013.08.12"

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
	log.Fatal(http.ListenAndServe(*address+":"+*port, nil)) //, Log(http.DefaultServeMux)))
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Server", fmt.Sprintf("goredirect/%s", VERSION))
	switch {
	// Perform a specific redirect response
	case rules[strings.Split(r.Host, ":")[0]][r.RequestURI] != nil:
		target = rules[strings.Split(r.Host, ":")[0]][r.RequestURI]
		redirectTarget(w, r, target)
		return
	// Perform a glob redirect response
	case rules[strings.Split(r.Host, ":")[0]]["*"] != nil:
		target = rules[strings.Split(r.Host, ":")[0]]["*"]
		// Parse the base url
		location, err := url.Parse(target["location"])
		if err != nil {
			internalServerError(w, r, target)
			return
		}
		// update with the request path, query
		location.Path = r.URL.Path
		location.RawQuery = r.URL.RawQuery
		// set the location value to the string location we built
		// from the r and the base location
		target["location"] = location.String()
		redirectTarget(w, r, target)
		return
	// Perform a default redirect response
	case rules[strings.Split(r.Host, ":")[0]]["default"] != nil:
		target = rules[strings.Split(r.Host, ":")[0]]["default"]
		redirectTarget(w, r, target)
		return
	}
	status = http.StatusNotFound
	http.NotFound(w, r)
	logit(r, NotFoundLength, status)
	return
}

func internalServerError(w http.ResponseWriter, r *http.Request, target map[string]string) {
	status = http.StatusInternalServerError
	setReviewed(w, target)
	message := fmt.Sprintf("<h4>Internal Server Error</h4><br />Unable to parse the response for <strong><em>%s</em></strong>", target["location"])
	size = len(message)
	http.Error(w, message, status)
	logit(r, size, status)
	return
}

func redirectTarget(w http.ResponseWriter, r *http.Request, target map[string]string) {
	status = http.StatusMovedPermanently
	setReviewed(w, target)
	size = len(target["location"]) + RedirectLength
	http.Redirect(w, r, target["location"], status)
	logit(r, size, status)
	return

}

func setReviewed(w http.ResponseWriter, target map[string]string) {
	if target["reviewed"] != "" {
		w.Header().Set("X-Last-Reviewed", target["reviewed"])
	}
}

func logit(r *http.Request, size int, status int) {
	t := time.Now()
	if r.URL.User != nil {
		user = r.URL.User.Username()
	} else {
		user = "-"
	}
	remoteAddr = r.Header.Get("X-Forwarded-For")
	if remoteAddr == "" {
		remoteAddr = strings.Split(r.RemoteAddr, ":")[0]
	}
	fmt.Printf("%s %s - %s [%s] \"%s %s %s\" %d %d \"%s\" \"%s\"\n",
		strings.Split(r.Host, ":")[0], remoteAddr, user, t.Format("02/Jan/2006:15:04:05 -0700"), r.Method, r.RequestURI, r.Proto, status, size, r.Referer(), r.UserAgent())
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
