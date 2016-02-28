package main

import (
	"io"
	"log"
	"net/http"
	"strings"
	//"math"
	"strconv"
)

// Make a shared map of state
// create a mutex around the shared map
// update key-value store for the improvisr

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

func handleOrientationRequest(w http.ResponseWriter, r *http.Request) {
	// parse the query args
	// id = <>, x = <> , y = <>
	// throws those values to another function to calculate and send values upstream
	query := r.URL.Query()
	log.Printf(r.URL.String())
	log.Printf("id = %s, x = %s, y = %s, c = %s", query["id"], query["x"], query["y"], query["c"])
	defer calculateVelocity(query["id"], query["x"], query["y"], query["c"])
	io.WriteString(w, "200")
}

func calculateVelocity(id, x, y, c []string) {
	log.Printf("Calculating velocity for: %s", id)
	sX := strings.Join(x, "")
	log.Printf(sX)
	iX, err := strconv.ParseInt(strings.Join(x, ""), 10, 32)
	iY, err := strconv.ParseInt(strings.Join(y, ""), 10, 32)
	iC, err := strconv.ParseInt(strings.Join(c, ""), 10, 32)
	if err != nil {
		log.Printf("Error occurred while calculating velocity")
	}
	log.Printf("x = %d", iX)
	log.Printf("y = %d", iY)
	log.Printf("c = %d", iC)
}

func main() {
	http.HandleFunc("/", hello)
	//http.HandleFunc("/arduino/", printIt)
	http.HandleFunc("/arduino/", proxyRequestToYun)
	http.HandleFunc("/orientation", handleOrientationRequest)
	http.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, r.URL.Path[1:])
	})
	panic(http.ListenAndServe(":8000", nil))
}
