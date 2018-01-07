package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// Town is the Node struct for each destination in the TSP
type Town struct {
	ID        int       `json:"id,omitEmpty"`
	Distances []float64 `json:"distances,omitEmpty"`
	Trails    []float64 `json:"-"`
	Score     float64   `json:"scores,omitEmpty"`
}

// Towns is the collection of nodes for the TSP, with a matrix of the probability of traversing between each town
type Towns struct {
	TownSlice         []Town
	ProbabilityMatrix [][]float64
}

// Town is the transfer object for converting Venue into an antcolony Town object
type TownTransferObject struct {
	ID        int       `json:"id,omitEmpty"`
	Distances []float64 `json:"distances,omitEmpty"`
	Score     float64   `json:"scores,omitEmpty"`
}

// Towns is the transfer object for converting Venues into an antcolony Towns object
type TownsTransferObject struct {
	TownSlice []TownTransferObject `json:"town"`
}

func createTownsFromDistances(ts []byte) (Towns, error) {
	// n := len(distances)
	//distances [][]float64, scores []float64
	//townsTransferObject := TownsTransferObject{}
	towns := Towns{}
	err := json.Unmarshal(ts, &towns)
	fmt.Println("ErrorUnmarshalling?:", err)
	//fmt.Println("tto: ", townsTransferObject)
	if err != nil {
		return towns, err
	}
	//n := len(townsTransferObject.TownSlice)
	n := len(towns.TownSlice)
	// towns := Towns{
	// 	TownSlice:         []Town{},
	// 	ProbabilityMatrix: make([][]float64, n),
	// }

	for i := range towns.ProbabilityMatrix {
		towns.ProbabilityMatrix[i] = make([]float64, n)
	}

	// for _, tto := range townsTransferObject.TownSlice {
	// 	town := Town{
	// 		ID:        tto.ID,
	// 		Distances: tto.Distances,
	// 		Score:     tto.Score,
	// 		Trails:    make([]float64, n),
	// 	}
	// 	for j := range town.Trails {
	// 		town.Trails[j] = 1
	// 	}
	// 	towns.TownSlice = append(towns.TownSlice, town)
	// }
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

func (t *Town) updateTrails(ants []ant) {
	// fmt.Println("Before:", (*t).trails)
	for i := range (*t).Trails {
		(*t).Trails[i] *= evaporationRate
	}
	// fmt.Println("Town:", (*t).id)
	for _, a := range ants {
		contribution := pheremoneStrength / a.score
		// a.printAnt()
		for j, myTour := range a.tour {
			if myTour == (*t).ID {
				// fmt.Println("Before:", (*t).trails)
				(*t).Trails[a.tour[j+1]] += contribution
				// fmt.Println("After:", (*t).trails)
				break
			}
		}
	}
	// fmt.Println("After:", (*t).trails)
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
