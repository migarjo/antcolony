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

var averageScore float64

func initializeGlobals() {
	numberOfTries = 50
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

func CreateTowns() Towns {
	if numberOfTries == 0 {
		initializeGlobals()
	}
	ts := createTowns(numberOfTowns, mapRange)
	return ts
}

func SolveTSP(ts Towns) ([]byte, []byte) {
	if numberOfTries == 0 {
		initializeGlobals()
	}

	var bestAnt ant
	var ants []ant

	ts.resetTrails()

	progressArray := ProgressOverTime{
		Generation:   []int{},
		AverageScore: []float64{},
	}
	for i := 0; i < numberOfTries; i++ {
		ants = createAntSlice(numberOfAnts, ts)

		if i > 0 {
			ants = append(ants, bestAnt)
		}

		for i := range ts.TownSlice {
			ts.TownSlice[i].updateTrails(ants)
			ts.clearProbabilityMatrix()
		}

		bestAnt, averageScore = analyzeAnts(ants)

		progressArray.add(averageScore, bestAnt.score)
	}
	so := createSigmaObject(&ts, &bestAnt)
	// fmt.Printf("%+v\n", so)

	soJSON := so.jsonify()
	progressArrayJSON := progressArray.jsonify()

	return soJSON, progressArrayJSON

}
