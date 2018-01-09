package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// Ant The representation of a single entity traversing a path throughout the towns
type Ant struct {
	ID            int
	Tour          []int
	Visited       []bool
	Probabilities []float64
	Score         float64
}

// ProgressOverTime tracks performance metrics of the ACO, such as AverageScore and MinimumScore for each Iteration
type ProgressOverTime struct {
	Iteration    []int     `json:"labels"`
	AverageScore []float64 `json:"average"`
	MinimumScore []float64 `json:"minimum"`
}

func createAnt(thisID int, townQty int) Ant {
	thisAnt := Ant{
		ID:            thisID,
		Tour:          []int{},
		Visited:       make([]bool, townQty),
		Probabilities: make([]float64, townQty),
		Score:         0,
	}
	return thisAnt
}

func (a *Ant) getProbabilityList(ts Towns) {
	// Current location
	i := (*a).Tour[len((*a).Tour)-1]
	n := len(ts.TownSlice)

	denom := 0.0
	numerator := make([]float64, n)

	for l := 0; l < n; l++ {
		if !(*a).Visited[l] {
			numerator[l] = ts.ProbabilityMatrix[i][l]
			denom += numerator[l]
		}
	}
	(*a).Probabilities[0] = 0
	for j := 1; j < n; j++ {
		if (*a).Visited[j] {
			(*a).Probabilities[j] = (*a).Probabilities[j-1]
		} else {
			(*a).Probabilities[j] = (*a).Probabilities[j-1] + numerator[j]/denom
		}
	}
	//fmt.Println("PM:", ts.probabilityMatrix)
	//fmt.Println("AntProbability", (*a).probabilities)
}

func (a *Ant) visitHome(ts Towns) {
	(*a).Tour = append((*a).Tour, ts.TownSlice[0].ID)
	(*a).Visited[ts.TownSlice[0].ID] = true
	if len((*a).Tour) > 1 {
		(*a).Score += ts.TownSlice[(*a).Tour[len((*a).Tour)-1]].Distances[(*a).Tour[len((*a).Tour)-2]]
	}
}

func (a *Ant) visitNextTown(ts Towns) {
	(*a).getProbabilityList(ts)
	randFloat := randSource.Float64()
	i := 0
	for randFloat > (*a).Probabilities[i] {
		i++
	}
	(*a).Tour = append((*a).Tour, i)
	(*a).Visited[i] = true
	(*a).Score += ts.TownSlice[i].Distances[(*a).Tour[len((*a).Tour)-2]]
}

func (a *Ant) printAnt() {
	fmt.Println((*a).Tour, (*a).Score)
}

func printAnts(a []Ant) {
	for _, Ant := range a {
		Ant.printAnt()
	}
}

func createAntSlice(n int, ts Towns) []Ant {
	ants := []Ant{}

	for a := 0; a < n; a++ {
		myAnt := createAnt(a, len(ts.TownSlice))

		myAnt.visitHome(ts)

		for len(myAnt.Tour) < len(ts.TownSlice) {
			myAnt.visitNextTown(ts)
		}

		myAnt.visitHome(ts)

		ants = append(ants, myAnt)
	}
	return ants
}

func analyzeAnts(ants []Ant) (Ant, float64) {
	bestAnt := ants[0]
	scoreTotal := 0.0
	for _, a := range ants {
		scoreTotal += a.Score
		if a.Score < bestAnt.Score {
			bestAnt = a
		}
	}
	return bestAnt, scoreTotal / float64(len(ants))
}

func (p *ProgressOverTime) add(averageScore float64, minimumScore float64) {
	(*p).Iteration = append((*p).Iteration, len((*p).Iteration))
	(*p).AverageScore = append((*p).AverageScore, averageScore)
	(*p).MinimumScore = append((*p).MinimumScore, minimumScore)
}

func (a *Ant) exportAnt() string {

	antJSON, err := json.Marshal(*a)

	if err != nil {
		fmt.Println(err)
	}
	return string(antJSON[:])
}

func (p *ProgressOverTime) jsonify() []byte {
	pJSON, err := json.Marshal(*p)

	if err != nil {
		fmt.Println(err)
	}
	return pJSON
}

func (p *ProgressOverTime) writeToFile(path string) {
	pJSON := (*p).jsonify()

	f, err := os.Create(path)
	if err != nil {
		fmt.Println(err)
	}

	n, err := f.Write(pJSON)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Wrote:", n, "objects")
}
