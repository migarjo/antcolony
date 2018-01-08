package main

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
)

// Town is the Node struct for each destination in the TSP
type Town struct {
	ID        int       `json:"id,omitEmpty"`
	Distances []float64 `json:"distances,omitEmpty"`
	Trails    []float64 `json:"trails,omitEmpty"`
	Score     float64   `json:"scores,omitEmpty"`
}

// Towns is the collection of nodes for the TSP, with a matrix of the probability of traversing between each town
type Towns struct {
	TownSlice         []Town      `json:"towns"`
	ProbabilityMatrix [][]float64 `json:"probabilitymatrix"`
}

func createTownsFromDistances(ts []byte) (Towns, error) {
	towns := Towns{}
	err := json.Unmarshal(ts, &towns)

	if err != nil {
		return towns, err
	}
	n := len(towns.TownSlice)

	if len(towns.TownSlice[0].Trails) == 0 {
		for i := range towns.TownSlice {
			towns.TownSlice[i].Trails = make([]float64, n)
		}
	}

	if len(towns.ProbabilityMatrix) == 0 {
		towns.ProbabilityMatrix = make([][]float64, n)
		for i := range towns.ProbabilityMatrix {
			towns.ProbabilityMatrix[i] = make([]float64, n)
		}
	}

	return towns, nil
}

func (ts *Towns) clearTrails() {
	for i := range (*ts).TownSlice {
		for j := range (*ts).TownSlice[i].Trails {
			(*ts).TownSlice[i].Trails[j] = 1
		}
	}
}

func (ts *Towns) clearProbabilityMatrix() {
	for i := range (*ts).ProbabilityMatrix {
		for j := range (*ts).ProbabilityMatrix[i] {
			(*ts).ProbabilityMatrix[i][j] = 0
		}
	}
}

func (t *Town) updateTrails(ants []Ant) {
	// fmt.Println("Before:", (*t).trails)
	for i := range (*t).Trails {
		(*t).Trails[i] *= evaporationRate
	}
	// fmt.Println("Town:", (*t).id)
	for _, a := range ants {
		contribution := pheremoneStrength / a.Score
		// a.printAnt()
		for j, myTour := range a.Tour {
			if myTour == (*t).ID {
				// fmt.Println("Before:", (*t).trails)
				(*t).Trails[a.Tour[j+1]] += contribution
				// fmt.Println("After:", (*t).trails)
				break
			}
		}
	}
	// fmt.Println("After:", (*t).trails)
}

func (ts *Towns) calculateProbabilityMatrix() {
	for i, t := range (*ts).TownSlice {
		for j := range (*ts).TownSlice[i].Trails {
			(*ts).ProbabilityMatrix[i][j] = math.Pow(t.Trails[j], trailPreference) * math.Pow(1.0/t.Distances[j], distancePreference)
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

func (ts *Towns) jsonify() []byte {
	tsJSON, err := json.Marshal(*ts)

	if err != nil {
		fmt.Println(err)
	}
	return tsJSON
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

func (ts *Towns) writeToFile(path string) {
	tsJSON := (*ts).jsonify()

	f, err := os.Create(path)
	if err != nil {
		fmt.Println(err)
	}

	n, err := f.Write(tsJSON)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Wrote:", n, "objects")
}
