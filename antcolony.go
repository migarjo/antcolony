package antcolony

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"
)

type model struct{}

var numberOfTries int
var numberOfTowns int
var numberOfAnts int
var mapRange int
var trailPreference float64
var distancePreference float64
var pheremoneStrength float64
var evaporationRate float64
var randSource *rand.Rand

var bestAnt ant
var averageScore float64

func initializeGlobals() {
	numberOfTries = 100
	numberOfTowns = 20
	numberOfAnts = 16
	mapRange = 20
	trailPreference = 1
	distancePreference = 1
	pheremoneStrength = 1
	evaporationRate = 0.8
	randSource = rand.New(rand.NewSource(time.Now().UnixNano()))
}

type message struct {
	Name string
	Body string
	Time int64
}

func OptimizeTSP() {
	initializeGlobals()

	towns := createTowns(numberOfTowns, mapRange)
	var bestAnt ant
	//towns := createBasicTowns(numberOfTowns)
	var ants []ant
	for i := 0; i < numberOfTries; i++ {
		ants = createAntSlice(numberOfAnts, towns)
		for i := range towns.townSlice {
			towns.townSlice[i].updateTrails(ants)
			towns.clearProbabilityMatrix()
		}
		//fmt.Printf("%+v", home)
		//fmt.Printf("%+v", towns)
		//fmt.Printf("%+v", ants)

		bestAnt, averageScore = analyzeAnts(ants)
		//fmt.Println("Best Ant:")
		//bestAnt.printAnt()
		fmt.Println("Average Score:", averageScore)
	}
	so := createSigmaObject(&towns, &bestAnt)
	fmt.Printf("%+v\n", so)

	soJSON, err := json.Marshal(so)

	if err != nil {
		fmt.Println(err)
	}

	f, err := os.Create("../simpleserver/data1.json")
	if err != nil {
		fmt.Println(err)
	}

	n, err := f.Write(soJSON)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Wrote:", n)

	printAnts(ants)
}
