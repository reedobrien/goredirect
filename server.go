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
)

var rules map[string]map[string]string

func main() {
	// log, err := syslog.New(syslog.LOG_ERR, "godir")
	// if err != nil {
	// 	panic("Could connect to syslog")
	// }
	// log.Info("Starting godir")
	address := flag.String("address", "127.0.0.1", "The address to listen on")
	path := flag.String("path", "", "Path to the csv of redirects")
	port := flag.String("port", "8080", "The port to listen on")
	// watch := flag.Bool("watch", false, "Watch for CSV file changes")
	flag.Parse()
	if *path == "" {
		log.Fatalln("You must supply a mapping file")
	}
	fmt.Println("Starting godir...")
	// rules := make(map[string]string)
	f, err := ioutil.ReadFile(*path)
	// defer f.Close()
	if err != nil {
		log.Panicln("Cant read file", err)
	}

	err = json.Unmarshal(f, &rules)
	if err != nil {
		log.Panicln("Can't read file", err)
	}

	http.HandleFunc("/", handler(redirectHandler, rules))
	log.Fatal(http.ListenAndServe(*address+":"+*port, nil))
	// reqpath := strings.Trim(req.URL.Path, "/")
}

func redirectHandler(w http.ResponseWriter, r *http.Request, rules map[string]map[string]string) {
	fmt.Println(rules)
}

func handler(fn func(http.ResponseWriter, *http.Request, map[string]map[string]string), rules map[string]map[string]string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		target := rules[strings.Split(r.Host, ":")[0]][r.URL.Path]
		fmt.Println("target", target)
		if target == "" {
			http.NotFound(w, r)
			return
		}
		http.Redirect(w, r, target, http.StatusMovedPermanently)
	}

	// fmt.Println(rules)
	// fmt.Fprintf(w, "Hi there, I love %s!\n", r.URL.Path)
}
