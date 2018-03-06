package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"reflect"
	"time"
)

var antRatio float64
var averageScore float64
var randSource *rand.Rand

func initializeGlobals() {
	antRatio = 0.8
	randSource = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func getStatusMessage() string {
	return "Our ants are ready to swarm!"
}

func status(w http.ResponseWriter, req *http.Request) {
	statusMessage := getStatusMessage()
	fmt.Fprintln(w, statusMessage)
}

func solvetsp(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Println("Error:", err)
		http.Error(w, fmt.Sprintf("Error %s", err),
			http.StatusInternalServerError)
		return
	}
	fmt.Println("Received body:")
	fmt.Println(string(body[:]))
	results, err := SolveTSP(body)
	if err != nil {
		if reflect.TypeOf(err) == reflect.TypeOf(ApplicationError{}) {
			log.Println("Error:", err)
			http.Error(w, fmt.Sprintf("Error: %s", err),
				http.StatusBadRequest)
			return
		}

		log.Println("Error:", err)
		http.Error(w, fmt.Sprintln("An internal server error has occurred. If problem persists, please contact support"),
			http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, results)
}

func main() {
	initializeGlobals()

	http.HandleFunc("/", status)
	http.HandleFunc("/status", status)
	http.HandleFunc("/api/solvetsp", solvetsp)
	var port string
	if os.Getenv("PORT") != "" {
		port = ":" + os.Getenv("PORT")
	} else {
		port = ":8000"
	}
	fmt.Println("Port set to", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Panicln(err)
	}
}
