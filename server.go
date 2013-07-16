package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

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
	rules := make(map[string]string)
	f, err := os.Open(*path)
	defer f.Close()
	if err != nil {
		panic("Cant read file")
	}

	reader := csv.NewReader(f)

	for {
		row, err := reader.Read()
		// stop at end of file
		if err == io.EOF {
			break
		}
		// ignore comments
		if strings.HasPrefix(row[0], "#") {
			continue
		}
		// build the rules map
		rules[row[0]] = row[1]
		// ignore empty lines
		if len(row) == 0 {
			break
		}
	}
	http.HandleFunc("/", handler(redirectHandler, rules))
	log.Fatal(http.ListenAndServe(*address+":"+*port, nil))
	// reqpath := strings.Trim(req.URL.Path, "/")
}

func redirectHandler(w http.ResponseWriter, r *http.Request, rules map[string]string) {
	fmt.Println(rules)
}

func handler(fn func(http.ResponseWriter, *http.Request, map[string]string), rules map[string]string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		target := rules[r.URL.Path]
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
