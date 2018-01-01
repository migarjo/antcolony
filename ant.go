package main

import "math"
import "fmt"

type ant struct {
	id            int
	tour          []int
	visited       []bool
	probabilities []float64
	score         float64
}

func createAnt(thisID int, townQty int) ant {
	thisAnt := ant{
		id:            thisID,
		tour:          []int{},
		visited:       make([]bool, townQty),
		probabilities: make([]float64, townQty),
		score:         0,
	}
	return thisAnt
}

func (a *ant) getProbabilityList(ts towns) {
	// Current location
	i := (*a).tour[len((*a).tour)-1]
	t := ts.townSlice[i]
	n := len(ts.townSlice)

	denom := 0.0
	numerator := make([]float64, n)

	for l := 0; l < n; l++ {
		if !(*a).visited[l] {
			if ts.probabilityMatrix[i][l] != 0 {
				numerator[l] = ts.probabilityMatrix[i][l]
				denom += numerator[l]
			} else {
				numerator[l] = math.Pow(t.trails[l], trailPreference) * math.Pow(1.0/t.distances[l], distancePreference)
				ts.probabilityMatrix[i][l] = numerator[l]
				denom += numerator[l]
			}
		}
	}
	(*a).probabilities[0] = 0
	for j := 1; j < n; j++ {
		if (*a).visited[j] {
			(*a).probabilities[j] = (*a).probabilities[j-1]
		} else {
			(*a).probabilities[j] = (*a).probabilities[j-1] + numerator[j]/denom
		}
	}
	//fmt.Println("PM:", ts.probabilityMatrix)
	//fmt.Println("AntProbability", (*a).probabilities)
}

func (a *ant) visitTown(t town, ts []town) {
	(*a).tour = append((*a).tour, t.id)
	(*a).visited[t.id] = true
	if len((*a).tour) > 1 {
		(*a).score += getDistance(t, ts[(*a).tour[len((*a).tour)-2]])
	}
}

func (a *ant) visitNextTown(ts towns) {
	(*a).getProbabilityList(ts)
	randFloat := randSource.Float64()
	i := 0
	for randFloat > (*a).probabilities[i] {
		i++
	}
	(*a).tour = append((*a).tour, i)
	(*a).visited[i] = true
	(*a).score += getDistance(ts.townSlice[i], ts.townSlice[(*a).tour[len((*a).tour)-2]])
}

func (a ant) printAnt() {
	fmt.Println(a.tour, a.score)
}

func printAnts(a []ant) {
	for _, ant := range a {
		ant.printAnt()
	}
}

func createAntSlice(n int, ts towns) []ant {
	ants := []ant{}

	for a := 0; a < n; a++ {
		myAnt := createAnt(a, len(ts.townSlice))
		myAnt.visitTown(ts.townSlice[0], ts.townSlice)

		for len(myAnt.tour) < len(ts.townSlice) {
			myAnt.visitNextTown(ts)
		}
		myAnt.visitTown(ts.townSlice[0], ts.townSlice)

		ants = append(ants, myAnt)
	}
	return ants
}

func analyzeAnts(ants []ant) (ant, float64) {
	bestAnt := ants[0]
	scoreTotal := 0.0
	for _, a := range ants {
		scoreTotal += a.score
		if a.score < bestAnt.score {
			bestAnt = a
		}
	}
	return bestAnt, scoreTotal / float64(len(ants))
}

func (a *ant) createPathFromTour() {

}
