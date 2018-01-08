package main

import (
	"encoding/json"
	"fmt"
)

type model struct{}

type Results struct {
	BestAnt       ant              `json:"bestant"`
	ProgressArray ProgressOverTime `json:"progress"`
	Towns         Towns            `json:"towns"`
}

func exportResults(a ant, p ProgressOverTime, ts Towns) string {
	results := Results{
		BestAnt:       a,
		ProgressArray: p,
		Towns:         ts,
	}

	resultsJSON, err := json.Marshal(results)

	if err != nil {
		fmt.Println(err)
	}
	return string(resultsJSON[:])
}

// SolveTSP for the Towns provided and returns JSON files for the optimum map and the average and best score over time
func SolveTSP(towns []byte) (string, error) {
	if numberOfTries == 0 {
		initializeGlobals()
	}

	var bestAnt ant
	var ants []ant

	ts, err := createTownsFromDistances(towns)
	if err != nil {
		return "", err
	}
	ts.clearTrails()

	progressArray := ProgressOverTime{
		Iteration:    []int{},
		AverageScore: []float64{},
		MinimumScore: []float64{},
	}
	for i := 0; i < numberOfTries; i++ {
		fmt.Println("i:", i)
		ants = createAntSlice(numberOfAnts, ts)

		if i > 0 {
			ants = append(ants, bestAnt)
		}

		for j := range ts.TownSlice {
			ts.TownSlice[j].updateTrails(ants)
			ts.clearProbabilityMatrix()
		}

		bestAnt, averageScore = analyzeAnts(ants)

		progressArray.add(averageScore, bestAnt.score)
	}
	//antJSON := exportAnt(bestAnt)
	// fmt.Printf("%+v\n", so)
	//fmt.Println(averageScore)
	//progressArrayJSON := string(progressArray.jsonify()[:])

	resultsJSON := exportResults(bestAnt, progressArray, ts)

	return resultsJSON, nil

}
