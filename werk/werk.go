package main

import (
	"io"
	"log"
	"net/http"
	"strings"
	//"math"
	//"gopkg.in/redis.v3"
	"strconv"
	"sync"
)

// Make a shared map of state
// create a mutex around the shared map
// update key-value store for the improvisr
var yRunningAverage = 0.0
var xRunningAverage = 0.0

var mutex = &sync.Mutex{}
var xs = make([]float64, 1)
var ys = make([]float64, 1)

var xAtTimes = make(map[int64][]float64)
var yAtTimes = make(map[int64][]float64)

func increaseCapacity(intSlice []float64) {
	newSlice := make([]float64, len(intSlice), (cap(intSlice)+1)*2)
	copy(newSlice, intSlice)
	intSlice = newSlice
}

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
		log.Printf("1: %g: %s", i, s[i])
		log.Printf("2: %g: %s", i, str)
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
	log.Printf("id = %s, x = %s, y = %s, c = %s, t = %s, ty = %s", query["id"], query["x"], query["y"], query["c"], query["t"], query["ty"])
	defer updateAverages(query["id"], query["x"], query["y"], query["c"], query["t"], query["ty"])
	io.WriteString(w, "200")
}

/*
func calculateVelocity(id, x, y, c []string) {
	log.Printf("Calculating velocity for: %s", id)
	sX := strings.Join(x, "")
	log.Printf(sX)
	iX, err := strconv.ParseFloat(strings.Join(x, ""), 10, 32)
	iY, err := strconv.ParseFloat(strings.Join(y, ""), 10, 32)
	iC, err := strconv.ParseFloat(strings.Join(c, ""), 10, 32)
	if err != nil {
		log.Printf("Error occurred while calculating velocity")
	}
	log.Printf("x = %g", iX)
	log.Printf("y = %g", iY)
	log.Printf("c = %g", iC)
}
*/

func updateAverages(id, x, y, c, t, ty []string) {
	log.Printf("Updating average: current id: %s", id)
	sX := strings.Join(x, "")
	log.Printf(sX)
	fX, err := strconv.ParseFloat(strings.Join(x, ""), 64)
	fY, err := strconv.ParseFloat(strings.Join(y, ""), 64)
	iC, err := strconv.ParseInt(strings.Join(c, ""), 10, 64)
	iT, err := strconv.ParseInt(strings.Join(t, ""), 10, 64)
	iTy, err := strconv.ParseInt(strings.Join(ty, ""), 10, 64)

	if err != nil {
		log.Printf("Error occurred while calculating velocity")
	}
	log.Printf("x = %g", fX)
	log.Printf("y = %g", fY)
	log.Printf("c = %d", iC)
	log.Printf("t = %d", iT)
	log.Printf("t = %d", iTy)
	//Lock around the xs
	mutex.Lock()
	//Check if we need to update the size of the arrays
	xs = xAtTimes[iT]
	xs = append(xs, fX)
	xAtTimes[iT] = xs
	//xs = xs[1:]
	log.Printf("xAtTimes[%d]: %v", iT, xAtTimes[iT])
	//ys = append(ys, fY)
	//ys = ys[1:]
	mutex.Unlock()
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
