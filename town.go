package antcolony

import (
	"math"
)

type town struct {
	id        int
	xCoord    int
	yCoord    int
	distances []float64
	trails    []float64
}

type towns struct {
	townSlice         []town
	probabilityMatrix [][]float64
}

func createBasicTowns(n int) towns {
	towns := towns{
		townSlice:         []town{},
		probabilityMatrix: make([][]float64, n),
	}
	for i := range towns.probabilityMatrix {
		towns.probabilityMatrix[i] = make([]float64, n)
	}

	for i := 0; i < n; i++ {
		thisTown := town{i, i, 0, []float64{}, make([]float64, n)}
		for i := range thisTown.trails {
			thisTown.trails[i] = 1
		}

		towns.townSlice = append(towns.townSlice, thisTown)
	}

	for i, ti := range towns.townSlice {
		//TODO: Test possible performance enhancement to directly write to townSlice(i) rather than ti
		for j, tj := range towns.townSlice {
			if i == j {
				ti.distances = append(ti.distances, 1)
			} else {
				ti.distances = append(ti.distances, getDistance(ti, tj))
			}
		}
		towns.townSlice[i] = ti
	}

	return towns
}

func createTowns(n int, fieldSize int) towns {

	towns := towns{
		townSlice:         []town{},
		probabilityMatrix: make([][]float64, n),
	}
	for i := range towns.probabilityMatrix {
		towns.probabilityMatrix[i] = make([]float64, n)
	}

	for i := 0; i < n; i++ {
		thisTown := town{i, randSource.Intn(fieldSize - 1), randSource.Intn(fieldSize - 1), []float64{}, make([]float64, n)}
		for i := range thisTown.trails {
			thisTown.trails[i] = 1
		}

		towns.townSlice = append(towns.townSlice, thisTown)
	}

	for i, ti := range towns.townSlice {
		//TODO: Test possible performance enhancement to directly write to townSlice(i) rather than ti
		for j, tj := range towns.townSlice {
			if i == j {
				ti.distances = append(ti.distances, 1)
			} else {
				ti.distances = append(ti.distances, getDistance(ti, tj))
			}
		}
		towns.townSlice[i] = ti
	}

	return towns
}

func (t *towns) clearProbabilityMatrix() {
	for i := range (*t).probabilityMatrix {
		for j := range (*t).probabilityMatrix[i] {
			(*t).probabilityMatrix[i][j] = 0
		}
	}
}

func (t *town) updateTrails(ants []ant) {
	// fmt.Println("Before:", (*t).trails)
	for i := range (*t).trails {
		(*t).trails[i] *= evaporationRate
	}
	// fmt.Println("Town:", (*t).id)
	for _, a := range ants {
		contribution := pheremoneStrength / a.score
		// a.printAnt()
		for j, myTour := range a.tour {
			if myTour == (*t).id {
				// fmt.Println("Before:", (*t).trails)
				(*t).trails[a.tour[j+1]] += contribution
				// fmt.Println("After:", (*t).trails)
				break
			}
		}
	}
	// fmt.Println("After:", (*t).trails)
}

func getDistance(ta town, tb town) float64 {
	return math.Sqrt(float64(ta.xCoord-tb.xCoord)*float64(ta.xCoord-tb.xCoord) + float64(ta.yCoord-tb.yCoord)*float64(ta.yCoord-tb.yCoord))
}
