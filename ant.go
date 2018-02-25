package main

import (
	"encoding/json"
	"fmt"
)

// Ant The representation of a single entity traversing a path throughout the towns
type Ant struct {
	ID               int       `json:"-"`
	Tour             []int     `json:"tour"`
	Visited          []bool    `json:"-"`
	Probabilities    []float64 `json:"-"`
	Score            float64   `json:"score"`
	DistanceTraveled float64   `json:"distanceTraveled"`
	AverageRating    float64   `json:"averageRating"`
}

// AverageResults tracks performance metrics of the ACO, such as AverageScore and MinimumScore for each Iteration
type AverageResults struct {
	Iteration int     `json:"labels"`
	Score     float64 `json:"score"`
	Distance  float64 `json:"distance"`
	Rating    float64 `json:"rating"`
}

func createAnt(thisID int, townQty int) Ant {
	thisAnt := Ant{
		ID:               thisID,
		Tour:             []int{},
		Visited:          make([]bool, townQty),
		Probabilities:    make([]float64, townQty),
		Score:            0,
		DistanceTraveled: 0,
		AverageRating:    0,
	}
	return thisAnt
}

func (a *Ant) getProbabilityList(ts Towns) {
	n := len(ts.TownSlice)
	denom := 0.0
	numerator := make([]float64, n)

	if len((*a).Tour) == 0 {
		for i := 0; i < n; i++ {
			if !(*a).Visited[i] {
				numerator[i] = ts.TownSlice[i].Rating
				denom += numerator[i]
			}
		}
	} else {
		currentLocation := (*a).Tour[len((*a).Tour)-1]
		for i := 0; i < n; i++ {
			if !(*a).Visited[i] {
				numerator[i] = ts.ProbabilityMatrix[currentLocation][i]
				denom += numerator[i]
			}
		}
	}

	if ts.IncludesHome {
		(*a).Probabilities[0] = 0
	} else {
		(*a).Probabilities[0] = numerator[0] / denom
	}

	for j := 1; j < n; j++ {
		if (*a).Visited[j] {
			(*a).Probabilities[j] = (*a).Probabilities[j-1]
		} else {
			(*a).Probabilities[j] = (*a).Probabilities[j-1] + numerator[j]/denom
		}
	}

}

func (a *Ant) visitHome(ts Towns) {
	if ts.IncludesHome {
		(*a).Tour = append((*a).Tour, ts.TownSlice[0].ID)
		(*a).Visited[ts.TownSlice[0].ID] = true
	}
}

func (a *Ant) returnHome(ts Towns) {
	if ts.IncludesHome {
		(*a).Tour = append((*a).Tour, ts.TownSlice[0].ID)
		distance := ts.TownSlice[(*a).Tour[len((*a).Tour)-1]].Distances[(*a).Tour[len((*a).Tour)-2]]
		(*a).Score += distance
		(*a).DistanceTraveled += distance
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
	tourLength := float64(len((*a).Tour))
	(*a).Visited[i] = true
	normalizedRating := ts.TownSlice[i].NormalizedRating
	if tourLength > 1 {
		distance := ts.TownSlice[i].Distances[(*a).Tour[len((*a).Tour)-2]]
		(*a).Score += 1 / (1/distance + normalizedRating)
		(*a).DistanceTraveled += distance
		(*a).AverageRating = (((*a).AverageRating*(tourLength-1) + ts.TownSlice[i].Rating) / tourLength)
	} else {
		if normalizedRating > 0 {
			(*a).Score += 1 / normalizedRating
		}

		(*a).AverageRating = ts.TownSlice[i].Rating
	}

}

func (a *Ant) printAnt() {
	fmt.Println((*a).Tour, (*a).Score)
}

func printAnts(a []Ant) {
	for _, Ant := range a {
		Ant.printAnt()
	}
}

func createAntSlice(n int, ts Towns, config AcoConfig) []Ant {
	ants := []Ant{}

	for a := 0; a < n; a++ {
		myAnt := createAnt(a, len(ts.TownSlice))

		myAnt.visitHome(ts)

		for len(myAnt.Tour) < config.VisitQuantity {
			myAnt.visitNextTown(ts)
		}

		myAnt.returnHome(ts)

		ants = append(ants, myAnt)
	}
	return ants
}

func analyzeAnts(ants []Ant, bestAnts []Ant, averageArray []AverageResults) ([]Ant, []AverageResults) {
	scoreTotal := 0.0
	distanceTotal := 0.0
	ratingTotal := 0.0
	bestAnt := ants[0]
	for _, a := range ants {
		scoreTotal += a.Score
		distanceTotal += a.DistanceTraveled
		ratingTotal += a.AverageRating
		if a.Score < bestAnt.Score {
			bestAnt = a
		}
	}
	averageResults := AverageResults{
		Iteration: len(averageArray),
		Score:     scoreTotal / float64(len(ants)),
		Distance:  distanceTotal / float64(len(ants)),
		Rating:    ratingTotal / float64(len(ants)),
	}

	if len(bestAnts) == 0 || bestAnt.Score < bestAnts[len(bestAnts)-1].Score {
		bestAnts = append(bestAnts, bestAnt)
	} else {
		bestAnts = append(bestAnts, bestAnts[len(bestAnts)-1])
	}

	averageArray = append(averageArray, averageResults)

	return bestAnts, averageArray
}

func (a *Ant) exportAnt() string {

	antJSON, err := json.Marshal(*a)

	if err != nil {
		fmt.Println(err)
	}
	return string(antJSON[:])
}
