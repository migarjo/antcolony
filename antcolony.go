package antcolony

import (
	"math/rand"
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

func SolveTSP() ([]byte, []byte) {
	initializeGlobals()

	towns := createTowns(numberOfTowns, mapRange)

	var bestAnt ant
	var ants []ant

	progressArray := ProgressOverTime{
		Generation:   []int{},
		AverageScore: []float64{},
	}
	for i := 0; i < numberOfTries; i++ {
		ants = createAntSlice(numberOfAnts, towns)

		if i > 0 {
			ants = append(ants, bestAnt)
		}

		for i := range towns.townSlice {
			towns.townSlice[i].updateTrails(ants)
			towns.clearProbabilityMatrix()
		}

		bestAnt, averageScore = analyzeAnts(ants)

		progressArray.add(averageScore, bestAnt.score)
	}
	so := createSigmaObject(&towns, &bestAnt)
	// fmt.Printf("%+v\n", so)

	soJSON := so.jsonify()
	progressArrayJSON := progressArray.jsonify()

	return soJSON, progressArrayJSON

}
