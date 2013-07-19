package main

import (
	"encoding/json"
	"flag"
	"fmt"
	// "io"
	"io/ioutil"
	"log"
	"net/http"
	// "os"
	"github.com/howeyc/fsnotify"
	"strings"
	"time"
)

var rules map[string]map[string]string
var size int = 19
var status int = http.StatusNotFound

func main() {
	address := flag.String("address", "127.0.0.1", "The address to listen on")
	path := flag.String("path", "", "Path to the json file of redirects")
	port := flag.String("port", "8080", "The port to listen on")
	// watch := flag.Bool("watch", false, "Watch for CSV file changes")
	flag.Parse()
	if *path == "" {
		log.Fatalln("You must supply a mapping file")
	}

	fmt.Printf("Starting godir...listening on http://%s:%s\n", *address, *port)

	log.Println("Loading rules from:", *path)
	err := loadRules(*path)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Starting watcher for ", *path)
	watcher, err := fsnotify.NewWatcher()
	watcher.Watch(*path)
	if err != nil {
		log.Println("Can't watch file", *path)
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
				log.Println("Error watching file:", *path, err)
			}
		}
	}()

	http.HandleFunc("/", handler(redirectHandler, rules))
	log.Fatal(http.ListenAndServe(*address+":"+*port, Log(http.DefaultServeMux)))
}

func redirectHandler(w http.ResponseWriter, r *http.Request, rules map[string]map[string]string) {
	// this function does nothing remove it and setup the handler to have
	// the right signature for http.HandleFunc
	log.Println(rules)
}

func handler(fn func(http.ResponseWriter, *http.Request, map[string]map[string]string), rules map[string]map[string]string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		target := rules[strings.Split(r.Host, ":")[0]][r.URL.Path]
		if target == "" {
			http.NotFound(w, r)
			return
		}
		status = http.StatusMovedPermanently
		size = 0
		http.Redirect(w, r, target, status)
		return
	}
}

func Log(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t := time.Now()
		fmt.Printf("%s - %s [%s] \"%s %s %s\" %d %d \"%s\" \"%s\"\n",
			strings.Split(r.RemoteAddr, ":")[0], r.URL.User, t.Format("02/Jan/2006:15:04:05 -0700"), r.Method, r.URL, r.Proto, status, size, r.Referer(), r.UserAgent())
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
