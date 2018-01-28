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
	ID              int       `json:"id,omitEmpty"`
	Distances       []float64 `json:"distances,omitEmpty"`
	Trails          []float64 `json:"trails,omitEmpty"`
	Score           float64   `json:"scores,omitEmpty"`
	NormalizedScore float64   `json:"-"`
}

// Towns is the collection of nodes for the TSP, with a matrix of the probability of traversing between each town
type Towns struct {
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

	(*ts).normalizeTownScores(config)

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
			if myTour == (*t).ID {
				(*t).Trails[a.Tour[j+1]] += contribution
				break
			}
		}
	}
}

func (ts *Towns) normalizeTownScores(config AcoConfig) {
	if config.ScorePreference == 0 {
		for i := range ts.TownSlice {
			(*ts).TownSlice[i].NormalizedScore = 0
		}
	} else {
		maxScore := ts.TownSlice[0].Score
		minScore := ts.TownSlice[0].Score
		minDistance := ts.TownSlice[0].Distances[1]
		for i := range ts.TownSlice {
			if ts.TownSlice[i].Score > maxScore {
				maxScore = ts.TownSlice[i].Score
			}
			if ts.TownSlice[i].Score < maxScore {
				minScore = ts.TownSlice[i].Score
			}
			for j := range ts.TownSlice[i].Distances {
				if ts.TownSlice[i].Distances[j] > 0 && ts.TownSlice[i].Distances[j] < minDistance {
					minDistance = ts.TownSlice[i].Distances[j]
				}
			}
		}
		maxDistanceFactor := 1.0 / minDistance

		if minScore == maxScore {
			for i := range ts.TownSlice {
				(*ts).TownSlice[i].NormalizedScore = 0
			}
		} else {
			if config.MaximizeScore == true {
				for i := range ts.TownSlice {
					(*ts).TownSlice[i].NormalizedScore = config.ScorePreference * maxDistanceFactor * ((*ts).TownSlice[i].Score - minScore) / (maxScore - minScore)
					fmt.Println("NormalizedScore in town:", i, "equals", (*ts).TownSlice[i].NormalizedScore)
				}
			}
		}
	}
}

func (ts *Towns) calculateProbabilityMatrix(config AcoConfig) {

	for i, t := range (*ts).TownSlice {
		for j := range (*ts).TownSlice[i].Trails {
			(*ts).ProbabilityMatrix[i][j] = math.Pow(t.Trails[j], config.TrailPreference) * math.Pow((1.0/t.Distances[j]+t.NormalizedScore), config.DistancePreference)
			fmt.Println(i, j, (*ts).ProbabilityMatrix[i][j])
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
