package main

import "math"
import "fmt"

type ant struct {
	tour          []int
	visited       []bool
	probabilities []float64
}

func createAnt(townQty int) ant {
	thisAnt := ant{
		tour:          []int{},
		visited:       make([]bool, townQty),
		probabilities: make([]float64, townQty),
	}
	return thisAnt
}

func (a *ant) getProbabilityList(ts []town) {
	// Current location
	i := (*a).tour[len((*a).tour)-1]
	t := ts[i]
	n := len(ts)

	denom := 0.0
	numerator := make([]float64, n)

	for l := 0; l < n; l++ {
		if !(*a).visited[l] {
			numerator[l] = math.Pow(t.trails[l], alpha) * math.Pow(1.0/t.distances[l], beta)
			denom += numerator[l]
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
}

func (a *ant) visitNextTown(ts []town) {
	(*a).getProbabilityList(ts)
	randFloat := randSource.Float64()
	i := 0
	fmt.Println(randFloat)
	fmt.Println((*a).probabilities[i])
	for randFloat > (*a).probabilities[i] {
		i++
	}
	(*a).tour = append((*a).tour, i)
	(*a).visited[i] = true
}

func (a *ant) visitTown(t town) {
	(*a).tour = append((*a).tour, t.id)
	(*a).visited[t.id] = true
}
