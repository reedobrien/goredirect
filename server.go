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
	"strings"
	"time"
)

var rules map[string]map[string]string
var status int = http.StatusNotFound
var size int = 19

func main() {
	address := flag.String("address", "127.0.0.1", "The address to listen on")
	path := flag.String("path", "", "Path to the csv of redirects")
	port := flag.String("port", "8080", "The port to listen on")
	// watch := flag.Bool("watch", false, "Watch for CSV file changes")
	flag.Parse()
	if *path == "" {
		log.Fatalln("You must supply a mapping file")
	}


	log.Println("Loading rules from:", *path)
	if err != nil {
		log.Panicln("Cant read file", err)
	}

	err = json.Unmarshal(f, &rules)
	if err != nil {
		log.Panicln("Can't read file", err)
	}

	http.HandleFunc("/", handler(redirectHandler, rules))
	log.Fatal(http.ListenAndServe(*address+":"+*port, Log(http.DefaultServeMux)))
}

func redirectHandler(w http.ResponseWriter, r *http.Request, rules map[string]map[string]string) {
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
