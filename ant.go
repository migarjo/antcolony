package main

import (
	"encoding/json"
	"fmt"
)

// Ant The representation of a single entity traversing a path throughout the towns
type Ant struct {
	ID            int       `json:"-"`
	Tour          []int     `json:"tour"`
	Visited       []bool    `json:"-"`
	Probabilities []float64 `json:"-"`
	Score         float64   `json:"score"`
	Distance      float64   `json:"distance"`
	Rating        float64   `json:"rating"`
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
		ID:            thisID,
		Tour:          []int{},
		Visited:       make([]bool, townQty),
		Probabilities: make([]float64, townQty),
		Score:         0,
		Distance:      0,
		Rating:        0,
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
		(*a).Distance += distance
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
		(*a).Distance += distance
		if ts.IncludesHome {
			(*a).Rating = ((*a).Rating*(tourLength-2) + ts.TownSlice[i].Rating) / (tourLength - 1)
		} else {
			(*a).Rating = (((*a).Rating*(tourLength-1) + ts.TownSlice[i].Rating) / tourLength)
		}
	} else {
		if normalizedRating > 0 {
			(*a).Score += 1 / normalizedRating
		}
		if !ts.IncludesHome {
			(*a).Rating = ts.TownSlice[i].Rating
		}
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
		requiredTownsVisited := make([]bool, len(ts.requiredTownsVisited))
		copy(requiredTownsVisited, ts.requiredTownsVisited)

		myAnt.visitHome(ts)

		if ts.IncludesHome {
			requiredTownsVisited[myAnt.Tour[len(myAnt.Tour)-1]] = true
		}

		for len(myAnt.Tour) < config.VisitQuantity {
			myAnt.visitNextTown(ts)
			requiredTownsVisited[myAnt.Tour[len(myAnt.Tour)-1]] = true
		}

		myAnt.returnHome(ts)

		unvisitedRequiredTowns := []int{}
		for i := range requiredTownsVisited {
			if !myAnt.Visited[i] && ts.TownSlice[i].IsRequired {
				unvisitedRequiredTowns = append(unvisitedRequiredTowns, i)
			}
			// if !visited {
			// 	unvisitedRequiredTowns = append(unvisitedRequiredTowns, i)
			// }
		}

		if len(unvisitedRequiredTowns) > 0 {
			fmt.Println("Tour Before: ", myAnt.Tour)
			myAnt.substituteRequiredTowns(ts, unvisitedRequiredTowns)
			fmt.Println("Tour After: ", myAnt.Tour)
		}

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
		distanceTotal += a.Distance
		ratingTotal += a.Rating
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

func (a *Ant) substituteRequiredTowns(ts Towns, unvisitedRequiredTowns []int) {

	for _, unvisitedTownID := range unvisitedRequiredTowns {
		a.substituteTown(ts, unvisitedTownID)
	}

	(*a).Distance = 0.0
	(*a).Score = 0.0
	ratingSum := 0.0
	for tourIndex, location := range (*a).Tour {
		normalizedRating := ts.TownSlice[tourIndex].NormalizedRating
		if tourIndex > 0 {
			distance := ts.TownSlice[location].Distances[(*a).Tour[tourIndex-1]]
			(*a).Score += 1 / (1/distance + normalizedRating)
			(*a).Distance += distance
			ratingSum += (*a).Rating
		} else {
			if normalizedRating > 0 {
				(*a).Score += 1 / normalizedRating
			}
			if !ts.IncludesHome {
				ratingSum += ts.TownSlice[location].Rating
			}
		}
	}
	if ts.IncludesHome {
		(*a).Rating = ratingSum / float64(len((*a).Tour)-2)
	} else {
		(*a).Rating = ratingSum / float64(len((*a).Tour))
	}

}

func (a *Ant) getReplaceProbabilityList(ts Towns, unvisitedTown int) {
	numeratorArray := make([]float64, len(ts.TownSlice))
	denom := 0.0

	for j, visitedTown := range a.Tour {
		if !ts.TownSlice[visitedTown].IsRequired {
			if j == 0 {
				numeratorArray[visitedTown] = 2 * ts.ProbabilityMatrix[unvisitedTown][2]
			} else if j == len(a.Tour)-1 {
				numeratorArray[visitedTown] = 2 * ts.ProbabilityMatrix[unvisitedTown][len(a.Tour)-2]
			} else {
				numeratorArray[visitedTown] = ts.ProbabilityMatrix[unvisitedTown][j-1] + ts.ProbabilityMatrix[unvisitedTown][j+1]
			}
			denom += numeratorArray[visitedTown]
		}
	}
	fmt.Println("Numerators:", numeratorArray, "denom:", denom)
	for i, numerator := range numeratorArray {
		if i == 0 {
			(*a).Probabilities[i] = numerator / denom
		} else {
			(*a).Probabilities[i] = (*a).Probabilities[i-1] + numerator/denom
		}
	}
}

func (a *Ant) substituteTown(ts Towns, unvisitedTownIndex int) {
	(*a).getReplaceProbabilityList(ts, unvisitedTownIndex)
	fmt.Println("Probability List: ", (*a).Probabilities)
	randFloat := randSource.Float64()
	replaceTownIndex := 0
	for randFloat > (*a).Probabilities[replaceTownIndex] {
		replaceTownIndex++
	}
	fmt.Println("Replace Index", replaceTownIndex, "Tour", (*a).Tour)
	replaceLocation := 0
	for replaceTownIndex != (*a).Tour[replaceLocation] {
		replaceLocation++
	}

	(*a).Visited[replaceTownIndex] = false
	(*a).Tour[replaceLocation] = unvisitedTownIndex
	(*a).Visited[unvisitedTownIndex] = true
}

func (a *Ant) exportAnt() string {

	antJSON, err := json.Marshal(*a)

	if err != nil {
		fmt.Println(err)
	}
	return string(antJSON[:])
}
