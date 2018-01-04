package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

var numberOfTries int
var numberOfAnts int
var trailPreference float64
var distancePreference float64
var pheremoneStrength float64
var evaporationRate float64
var averageScore float64
var randSource *rand.Rand

func initializeGlobals() {
	numberOfTries = 50
	numberOfAnts = 16
	trailPreference = 1
	distancePreference = 1
	pheremoneStrength = 1
	evaporationRate = 0.8
	randSource = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func status(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(w, "Our ants are ready to swarm!")
}

func solvetsp(w http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}
	log.Println(string(body))
	solution, progressArray, err := SolveTSP(body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error %s", err),
			http.StatusInternalServerError)
		return
	}
	log.Println(progressArray)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	fmt.Fprintln(w, solution)
}

func main() {
	initializeGlobals()

	http.HandleFunc("/status", status)
	http.HandleFunc("/api/solvetsp", solvetsp)
	var port string
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	} else {
		port = ":8000"
	}
	http.ListenAndServe(port, nil)
}
