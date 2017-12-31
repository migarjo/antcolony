package main

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

func createTowns(n int, fieldSize int) []town {

	townSlice := []town{}
	for i := 0; i < n; i++ {
		thisTown := town{i, randSource.Intn(fieldSize - 1), randSource.Intn(fieldSize - 1), []float64{}, make([]float64, n)}
		for i := range thisTown.trails {
			thisTown.trails[i] = 1
		}

		townSlice = append(townSlice, thisTown)
	}

	for i, ti := range townSlice {
		//TODO: Test possible performance enhancement to directly write to townSlice(i) rather than ti
		for j, tj := range townSlice {
			if i == j {
				ti.distances = append(ti.distances, 1)
			} else {
				ti.distances = append(ti.distances, getDistance(ti, tj))
			}
		}
		townSlice[i] = ti
	}

	return townSlice
}

func getDistance(ta town, tb town) float64 {
	return math.Sqrt(float64(ta.xCoord-tb.xCoord)*float64(ta.xCoord-tb.xCoord) + float64(ta.yCoord-tb.yCoord)*float64(ta.yCoord-tb.yCoord))
}
