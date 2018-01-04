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
	XCoord    int       `json:"xCoord,omitEmpty"`
	YCoord    int       `json:"yCoord,omitEmpty"`
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
	townsTransferObject := TownsTransferObject{}
	err := json.Unmarshal(ts, &townsTransferObject)
	fmt.Println("ErrorUnmarshalling?:", err)
	fmt.Println("tto: ", townsTransferObject)
	if err != nil {
		return Towns{}, err
	}
	n := len(townsTransferObject.TownSlice)
	towns := Towns{
		TownSlice:         []Town{},
		ProbabilityMatrix: make([][]float64, n),
	}

	for i := range towns.ProbabilityMatrix {
		towns.ProbabilityMatrix[i] = make([]float64, n)
	}

	for _, tto := range townsTransferObject.TownSlice {
		town := Town{
			ID:        tto.ID,
			Distances: tto.Distances,
			Score:     tto.Score,
			Trails:    make([]float64, n),
		}
		for j := range town.Trails {
			town.Trails[j] = 1
		}
		towns.TownSlice = append(towns.TownSlice, town)
	}
	return towns, nil
}

func createBasicTowns(n int) Towns {
	towns := Towns{
		TownSlice:         []Town{},
		ProbabilityMatrix: make([][]float64, n),
	}
	for i := range towns.ProbabilityMatrix {
		towns.ProbabilityMatrix[i] = make([]float64, n)
	}

	for i := 0; i < n; i++ {
		thisTown := Town{
			ID:        i,
			XCoord:    i,
			YCoord:    0,
			Distances: []float64{},
			Trails:    make([]float64, n),
			Score:     1,
		}
		for i := range thisTown.Trails {
			thisTown.Trails[i] = 1
		}

		towns.TownSlice = append(towns.TownSlice, thisTown)
	}

	for i, ti := range towns.TownSlice {
		//TODO: Test possible performance enhancement to directly write to townSlice(i) rather than ti
		for j, tj := range towns.TownSlice {
			if i == j {
				ti.Distances = append(ti.Distances, 1)
			} else {
				ti.Distances = append(ti.Distances, getDistanceFromXY(ti, tj))
			}
		}
		towns.TownSlice[i] = ti
	}

	return towns
}

func createTowns(n int, fieldSize int) Towns {

	towns := Towns{
		TownSlice:         []Town{},
		ProbabilityMatrix: make([][]float64, n),
	}
	for i := range towns.ProbabilityMatrix {
		towns.ProbabilityMatrix[i] = make([]float64, n)
	}

	for i := 0; i < n; i++ {
		thisTown := Town{
			ID:        i,
			XCoord:    randSource.Intn(fieldSize - 1),
			YCoord:    randSource.Intn(fieldSize - 1),
			Distances: []float64{},
			Trails:    make([]float64, n),
			Score:     1,
		}
		for i := range thisTown.Trails {
			thisTown.Trails[i] = 1
		}

		towns.TownSlice = append(towns.TownSlice, thisTown)
	}

	for i, ti := range towns.TownSlice {
		//TODO: Test possible performance enhancement to directly write to townSlice(i) rather than ti
		for j, tj := range towns.TownSlice {
			if i == j {
				ti.Distances = append(ti.Distances, 1)
			} else {
				ti.Distances = append(ti.Distances, getDistanceFromXY(ti, tj))
			}
		}
		towns.TownSlice[i] = ti
	}

	return towns
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

func getDistanceFromXY(ta Town, tb Town) float64 {
	return math.Sqrt(float64(ta.XCoord-tb.XCoord)*float64(ta.XCoord-tb.XCoord) + float64(ta.YCoord-tb.YCoord)*float64(ta.YCoord-tb.YCoord))
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
