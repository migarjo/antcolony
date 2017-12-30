package main

import (
	"math/rand"
	"time"
)

type town struct {
	id     int
	xCoord int
	yCoord int
}

func createTowns(n int, fieldSize int) []town {
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)

	townSlice := []town{}
	for i := 0; i < n; i++ {
		r.Intn(fieldSize - 1)
		thisTown := town{i, r.Intn(fieldSize - 1), r.Intn(fieldSize - 1)}
		townSlice = append(townSlice, thisTown)
	}
	return townSlice
}
