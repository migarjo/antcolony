package main

import (
	"fmt"
	"math/rand"
	"time"
)

type model struct{}

var alpha float64
var beta float64
var randSource *rand.Rand

func main() {
	alpha = 1
	beta = 5
	randSource = rand.New(rand.NewSource(time.Now().UnixNano()))
	towns := createTowns(10, 100)

	myAnt := createAnt(10)
	myAnt.visitTown(towns[0])

	for len(myAnt.tour) < 10 {
		fmt.Println(myAnt.tour)
		myAnt.visitNextTown(towns)
	}

	//fmt.Printf("%+v", home)
	//fmt.Printf("%+v", towns)
	fmt.Printf("%+v", myAnt)
}
