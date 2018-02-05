package main

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"strconv"
)

// Town is the Node struct for each destination in the TSP
type Town struct {
	ID               int       `json:"id,omitEmpty"`
	Distances        []float64 `json:"distances,omitEmpty"`
	Trails           []float64 `json:"trails,omitEmpty"`
	Rating           float64   `json:"rating,omitEmpty"`
	IsRequired       bool      `json:"isRequired"`
	NormalizedRating float64   `json:"-"`
}

// Towns is the collection of nodes for the TSP, with a matrix of the probability of traversing between each town
type Towns struct {
	IncludesHome      bool        `json:"includesHome"`
	TownSlice         []Town      `json:"towns"`
	ProbabilityMatrix [][]float64 `json:"-"`
}

func (ts *Towns) initializeTowns(config AcoConfig) error {
	n := len((*ts).TownSlice)

	if n == 0 {
		return ApplicationError{"No towns provided"}
	}
	if n == 1 {
		return ApplicationError{"Only one town provided"}
	}

	for i, t := range (*ts).TownSlice {
		if len(t.Distances) != len((*ts).TownSlice) {
			return ApplicationError{"Number of distances for town: " + strconv.Itoa(t.ID) + " is inconsistent with total number of towns"}
		}
		if len(t.Trails) == 0 {
			(*ts).TownSlice[i].Trails = make([]float64, n)
			for j := range (*ts).TownSlice[i].Trails {
				(*ts).TownSlice[i].Trails[j] = 1
			}
		}
	}

	if len((*ts).ProbabilityMatrix) == 0 {
		(*ts).ProbabilityMatrix = make([][]float64, n)
		for i := range (*ts).ProbabilityMatrix {
			(*ts).ProbabilityMatrix[i] = make([]float64, n)
		}
	}

	(*ts).normalizeTownRatings(config)

	return nil
}

func (ts *Towns) clearProbabilityMatrix() {
	for i := range (*ts).ProbabilityMatrix {
		for j := range (*ts).ProbabilityMatrix[i] {
			(*ts).ProbabilityMatrix[i][j] = 0
		}
	}
}

func (t *Town) updateTrails(ants []Ant, config AcoConfig) {
	for i := range (*t).Trails {
		(*t).Trails[i] *= config.EvaporationRate
	}
	for _, a := range ants {
		contribution := config.PheremoneStrength / a.Score
		for j, myTour := range a.Tour {
			if myTour == (*t).ID && j != len(a.Tour)-1 {
				(*t).Trails[a.Tour[j+1]] += contribution
				break
			}
		}
	}
}

func (ts *Towns) normalizeTownRatings(config AcoConfig) {
	if config.RatingPreference == 0 {
		for i := range ts.TownSlice {
			(*ts).TownSlice[i].NormalizedRating = 1
		}
	} else {
		maxRating := ts.TownSlice[0].Rating
		minRating := ts.TownSlice[0].Rating
		minDistance := ts.TownSlice[0].Distances[1]
		for i := range ts.TownSlice {
			if ts.TownSlice[i].Rating > maxRating {
				maxRating = ts.TownSlice[i].Rating
			}
			if ts.TownSlice[i].Rating < maxRating {
				minRating = ts.TownSlice[i].Rating
			}
			for j := range ts.TownSlice[i].Distances {
				if ts.TownSlice[i].Distances[j] > 0 && ts.TownSlice[i].Distances[j] < minDistance {
					minDistance = ts.TownSlice[i].Distances[j]
				}
			}
		}
		maxDistanceFactor := 1.0 / minDistance
		if minRating == maxRating {
			for i := range ts.TownSlice {
				(*ts).TownSlice[i].NormalizedRating = 1
			}
		} else {
			if config.MaximizeRating == true {
				for i := range ts.TownSlice {
					(*ts).TownSlice[i].NormalizedRating = config.RatingPreference * maxDistanceFactor * ((*ts).TownSlice[i].Rating - minRating) / (maxRating - minRating)
					fmt.Println("NormalizedRating in town:", i, "equals", (*ts).TownSlice[i].NormalizedRating)
				}
			}
		}
	}
}

func (ts *Towns) calculateProbabilityMatrix(config AcoConfig) {

	for i, t := range (*ts).TownSlice {
		for j := range (*ts).TownSlice[i].Trails {
			(*ts).ProbabilityMatrix[i][j] = math.Pow(t.Trails[j], config.TrailPreference) * math.Pow((1.0/t.Distances[j]+t.NormalizedRating), config.DistancePreference)
		}
	}
}

func (t *Town) jsonify() []byte {
	tJSON, err := json.Marshal(*t)

	if err != nil {
		fmt.Println(err)
	}
	return tJSON
}

func (t *Town) writeToFile(path string) {
	tJSON := (*t).jsonify()

	f, err := os.Create(path)
	if err != nil {
		fmt.Println(err)
	}

	n, err := f.Write(tJSON)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Wrote:", n, "objects")
}
