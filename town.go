package antcolony

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
)

// Town is the Node struct for each destination in the TSP
type Town struct {
	ID         int
	XCoord     int
	YCoord     int
	Distances  []float64
	Trails     []float64
	Properties map[string]string
}

// Towns is the collection of nodes for the TSP, with a matrix of the probability of traversing between each town
type Towns struct {
	TownSlice         []Town
	ProbabilityMatrix [][]float64
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
			ID:         i,
			XCoord:     i,
			YCoord:     0,
			Distances:  []float64{},
			Trails:     make([]float64, n),
			Properties: make(map[string]string),
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
				ti.Distances = append(ti.Distances, getDistance(ti, tj))
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
				ti.Distances = append(ti.Distances, getDistance(ti, tj))
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

func getDistance(ta Town, tb Town) float64 {
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
