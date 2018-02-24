package main

import (
	"encoding/json"
	"fmt"
	"math"
)

//AcoConfig Configuration parameters for the ACO algorithm
type AcoConfig struct {
	NumberOfIterations int     `json:"numberOfIterations"`
	TrailPreference    float64 `json:"trailPreference"`
	RatingPreference   float64 `json:"ratingPreference"`
	DistancePreference float64 `json:"distancePreference"`
	PheremoneStrength  float64 `json:"pheremoneStrength"`
	EvaporationRate    float64 `json:"evaporationRate"`
	MaximizeRating     bool    `json:"maximizeRating"`
	VisitQuantity      int     `json:"visitQuantity"`
	Verbose            bool    `json:"verbose"`
}

// Inputs Input parameters, including configuration and towns
type Inputs struct {
	AcoConfig `json:"config"`
	Towns     `json:"towns"`
}

// Results The best ant, progress array, and towns to be returned from the web service
type Results struct {
	BestAnts      []Ant            `json:"bestAnts"`
	ProgressArray ProgressOverTime `json:"progress"`
	Towns         Towns            `json:"towns"`
}

func exportResults(as []Ant, p ProgressOverTime, ts Towns) (string, error) {
	results := Results{
		BestAnts:      as,
		ProgressArray: p,
		Towns:         ts,
	}

	resultsJSON, err := json.Marshal(results)

	if err != nil {
		fmt.Println("Error serializing JSON:", err)
		fmt.Println(results)
		return "", err
	}
	return string(resultsJSON[:]), nil
}

func importInputs(inputsJSON []byte) (AcoConfig, Towns, error) {
	inputs := Inputs{
		AcoConfig: AcoConfig{
			NumberOfIterations: 50,
			TrailPreference:    1,
			DistancePreference: 1,
			RatingPreference:   0,
			MaximizeRating:     true,
			PheremoneStrength:  1,
			EvaporationRate:    .8,
			VisitQuantity:      0,
			Verbose:            false,
		},
		Towns: Towns{
			IncludesHome: true,
		},
	}

	err := json.Unmarshal(inputsJSON, &inputs)
	if err != nil {
		return inputs.AcoConfig, inputs.Towns, ApplicationError{"Error parsing input JSON: " + err.Error()}
	}
	fmt.Println(inputs.AcoConfig.PheremoneStrength)
	if inputs.AcoConfig.VisitQuantity == 0 {
		inputs.AcoConfig.VisitQuantity = len(inputs.Towns.TownSlice)
	}

	return inputs.AcoConfig, inputs.Towns, nil
}

// SolveTSP for the Towns provided and returns JSON files for the optimum map and the average and best score over time
func SolveTSP(towns []byte) (string, error) {

	config, ts, err := importInputs(towns)
	if err != nil {
		return "", err
	}

	var bestAnts []Ant
	var ants []Ant

	err = ts.initializeTowns(config)
	if err != nil {
		return "", err
	}

	numberOfAnts := int(math.Ceil(antRatio * float64(len(ts.TownSlice))))

	progressArray := ProgressOverTime{
		Iteration:    []int{},
		AverageScore: []float64{},
		MinimumScore: []float64{},
	}
	for i := 0; i < config.NumberOfIterations; i++ {
		// fmt.Println("Iteration: ", i)
		// for _, t := range ts.TownSlice {
		// 	fmt.Println(t.Trails)
		// }

		ts.calculateProbabilityMatrix(config)
		ants = createAntSlice(numberOfAnts, ts, config)

		if i > 0 {
			ants = append(ants, bestAnts[len(bestAnts)-1])
		}

		for j := range ts.TownSlice {
			ts.TownSlice[j].updateTrails(ants, config)
		}

		bestAnts, averageScore = analyzeAnts(ants, bestAnts)

		progressArray.add(averageScore, bestAnts[len(bestAnts)-1].Score)
	}

	resultsJSON, err := exportResults(bestAnts, progressArray, ts)
	if err != nil {
		return "", err
	}

	return resultsJSON, nil

}
