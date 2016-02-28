package main

import (
	"io"
	"log"
	"net/http"
	"strings"
)

func hello(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello, world!")
}

func printIt(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
}

func proxyRequestToYun(w http.ResponseWriter, r *http.Request) {
	log.Printf("Serving %s", r.URL)
	s := []string{"http://arduino.local", r.URL.Path}
	for i, str := range s {
		log.Printf("1: %d: %s", i, s[i])
		log.Printf("2: %d: %s", i, str)
	}
	st := strings.Join(s, "")
	res, err := http.Get(st)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Proxied request to arduino: %s", st)
	io.WriteString(w, res.Status)
}

func main() {
	http.HandleFunc("/", hello)
	//http.HandleFunc("/arduino/", printIt)
	http.HandleFunc("/arduino/", proxyRequestToYun)
	http.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, r.URL.Path[1:])
	})
	panic(http.ListenAndServe(":8000", nil))
}
