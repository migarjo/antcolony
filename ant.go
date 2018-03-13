package main

import (
	"encoding/json"
	"fmt"
	"math"
)

// Ant The representation of a single entity traversing a path throughout the towns
type Ant struct {
	ID                 int         `json:"-"`
	Tour               []int       `json:"tour"`
	Visited            []bool      `json:"-"`
	Probabilities      []float64   `json:"-"`
	Score              float64     `json:"score"`
	Distance           float64     `json:"distance"`
	Rating             float64     `json:"rating"`
	VisitSpan          [][]float64 `json:"visitSpan"`
	costState          float64
	availabilityBounds AvailabilityBounds
	tourComplete       bool
}

// AverageResults tracks performance metrics of the ACO, such as AverageScore and MinimumScore for each Iteration
type AverageResults struct {
	Iteration int     `json:"labels"`
	Score     float64 `json:"score"`
	Distance  float64 `json:"distance"`
	Rating    float64 `json:"rating"`
}

func createAnt(thisID int, townQty int, config AcoConfig) Ant {
	thisAnt := Ant{
		ID:            thisID,
		Tour:          []int{},
		Visited:       make([]bool, townQty),
		Probabilities: make([]float64, townQty),
		Score:         0,
		Distance:      0,
		Rating:        0,
		VisitSpan:     [][]float64{},
		availabilityBounds: AvailabilityBounds{
			Start: config.AvailabilityBounds.Start,
			End:   config.AvailabilityBounds.End,
		},
	}
	return thisAnt
}

func (a *Ant) getProbabilityList(ts Towns) {
	n := len(ts.TownSlice)
	denom := 0.0
	numerator := make([]float64, n)

	if len((*a).Tour) == 0 {
		for i := 0; i < n; i++ {
			if !(*a).Visited[i] && isAvailable(ts, a, i) {
				numerator[i] = ts.TownSlice[i].Rating
				denom += numerator[i]
			}
		}
		if denom == 0 {
			for i := 0; i < n; i++ {
				if !(*a).Visited[i] && isAvailable(ts, a, i) {
					numerator[i] = 1
					denom++
				}
			}
		}
	} else {
		currentLocation := (*a).Tour[len((*a).Tour)-1]
		for i := 0; i < n; i++ {
			if !(*a).Visited[i] && isAvailable(ts, a, i) {
				numerator[i] = ts.probabilityMatrix[currentLocation][i]
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
		(*a).VisitSpan = append((*a).VisitSpan, []float64{ts.AvailabilityBounds.Start, ts.AvailabilityBounds.Start + ts.TownSlice[0].VisitDuration})
		(*a).costState += ts.TownSlice[0].VisitDuration
	}
}

func (a *Ant) returnHome(ts Towns) {
	if ts.IncludesHome {
		(*a).Tour = append((*a).Tour, ts.TownSlice[0].ID)
		distance := ts.TownSlice[(*a).Tour[len((*a).Tour)-1]].Distances[(*a).Tour[len((*a).Tour)-2]]
		(*a).Score += distance
		(*a).Distance += distance
		thisVisitSpan := make([]float64, 2)
		thisVisitSpan[0] = (*a).VisitSpan[len((*a).VisitSpan)-1][1] + distance
		thisVisitSpan[1] = thisVisitSpan[0] + ts.TownSlice[0].VisitDuration
		(*a).VisitSpan = append((*a).VisitSpan, thisVisitSpan)
	}
}

func (a *Ant) visitNextTown(ts Towns) {
	randFloat := randSource.Float64()
	i := 0
	for randFloat > (*a).Probabilities[i] {
		i++
	}
	(*a).Tour = append((*a).Tour, i)
	(*a).Visited[i] = true
	normalizedRating := ts.TownSlice[i].NormalizedRating
	if tourLength > 1 {
		distance := ts.TownSlice[(*a).Tour[len((*a).Tour)-2]].Distances[i]
		(*a).Score += 1 / (1/distance + normalizedRating)
		(*a).Distance += distance
		if ts.IncludesHome {
			(*a).Rating = ((*a).Rating*(tourLength-2) + ts.TownSlice[i].Rating) / (tourLength - 1)
		} else {
			(*a).Rating = (((*a).Rating*(tourLength-1) + ts.TownSlice[i].Rating) / tourLength)
		}
		(*a).costState += ts.TownSlice[i].VisitDuration + distance
		thisVisitSpan := make([]float64, 2)
		thisVisitSpan[0] = (*a).VisitSpan[len((*a).VisitSpan)-1][1] + distance
		thisVisitSpan[1] = thisVisitSpan[0] + ts.TownSlice[i].VisitDuration
		(*a).VisitSpan = append((*a).VisitSpan, thisVisitSpan)
	} else {
		if normalizedRating > 0 {
			(*a).Score += 1 / normalizedRating
		}
		if !ts.IncludesHome {
			(*a).Rating = ts.TownSlice[i].Rating
		}
		thisVisitSpan := make([]float64, 2)
		thisVisitSpan[0] = ts.AvailabilityBounds.Start
		thisVisitSpan[1] = thisVisitSpan[0] + ts.TownSlice[i].VisitDuration
		(*a).VisitSpan = append((*a).VisitSpan, thisVisitSpan)
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

func (a *Ant) isTourComplete(ts Towns, config AcoConfig) bool {
	if len((*a).Tour) >= config.VisitQuantity {
		(*a).tourComplete = true
		return true
	}

	(*a).getProbabilityList(ts)
	fmt.Println((*a).VisitSpan, (*a).Probabilities)
	if math.IsNaN((*a).Probabilities[0]) {
		(*a).tourComplete = true
		return true
	}
	return false
}

func createAntSlice(n int, ts Towns, config AcoConfig) []Ant {
	ants := []Ant{}

	for a := 0; a < n; a++ {
		myAnt := Ant{}
		isTourLongEnough := false
		for !isTourLongEnough {
			tryCt := 0
			for !isTourLongEnough && tryCt < 100 {
				myAnt = createAnt(a, len(ts.TownSlice), config)

				myAnt.visitHome(ts)

				for !myAnt.isTourComplete(ts, config) {
					myAnt.visitNextTown(ts)
				}
				isTourLongEnough = getTourCapacityRatio(myAnt) > config.MinimumTripUsage
				tryCt++
			}
			config.MinimumTripUsage *= .9
		}
		myAnt.returnHome(ts)
		ants = append(ants, myAnt)
	}

	return ants
}

func getTourCapacityRatio(ant Ant) float64 {
	return (ant.VisitSpan[len(ant.VisitSpan)-1][1] - ant.availabilityBounds.Start) / (ant.availabilityBounds.End - ant.availabilityBounds.Start)
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

	randSource.Shuffle(len(unvisitedRequiredTowns), func(i, j int) {
		unvisitedRequiredTowns[i], unvisitedRequiredTowns[j] = unvisitedRequiredTowns[j], unvisitedRequiredTowns[i]
	})

	for _, unvisitedTownID := range unvisitedRequiredTowns {
		a.substituteTown(ts, unvisitedTownID)
	}
	a.Score, a.Distance, a.Rating = calculateAntResults((*a).Tour, ts)
}

func calculateAntResults(tour []int, ts Towns) (float64, float64, float64) {
	distance := 0.0
	score := 0.0
	rating := 0.0
	ratingSum := 0.0
	for tourIndex, location := range tour {
		normalizedRating := ts.TownSlice[location].NormalizedRating
		if tourIndex > 0 {
			legDistance := ts.TownSlice[location].Distances[tour[tourIndex-1]]
			score += 1 / (1/legDistance + normalizedRating)
			distance += legDistance
			ratingSum += ts.TownSlice[location].Rating
		} else {
			if normalizedRating > 0 {
				score += 1 / normalizedRating
			}
			if !ts.IncludesHome {
				ratingSum += ts.TownSlice[location].Rating
			}
		}
	}
	if ts.IncludesHome {
		rating = ratingSum / float64(len(tour)-2)
	} else {
		rating = ratingSum / float64(len(tour))
	}
	return score, distance, rating
}

func (a *Ant) getReplaceProbabilityList(ts Towns, unvisitedTown int) {
	originalNumeratorArray := make([]float64, len(ts.TownSlice))
	replacementNumeratorArray := make([]float64, len(ts.TownSlice))
	numeratorDiffArray := make([]float64, len(ts.TownSlice))
	minimumNumeratorDiff := 0.0
	diffDenom := 0.0

	for j, visitedTown := range a.Tour {
		if !ts.TownSlice[visitedTown].IsRequired {
			if j == 0 {
				originalNumeratorArray[visitedTown] = 2 * ts.probabilityMatrix[visitedTown][a.Tour[1]]
				replacementNumeratorArray[visitedTown] = 2 * ts.probabilityMatrix[unvisitedTown][a.Tour[1]]
			} else if j == len(a.Tour)-1 {
				originalNumeratorArray[visitedTown] = 2 * ts.probabilityMatrix[visitedTown][a.Tour[len(a.Tour)-2]]
				replacementNumeratorArray[visitedTown] = 2 * ts.probabilityMatrix[unvisitedTown][a.Tour[len(a.Tour)-2]]
			} else {
				originalNumeratorArray[visitedTown] = ts.probabilityMatrix[visitedTown][a.Tour[j-1]] + ts.probabilityMatrix[visitedTown][a.Tour[j+1]]
				replacementNumeratorArray[visitedTown] = ts.probabilityMatrix[unvisitedTown][a.Tour[j-1]] + ts.probabilityMatrix[unvisitedTown][a.Tour[j+1]]
			}
			numeratorDiffArray[visitedTown] = replacementNumeratorArray[visitedTown] - originalNumeratorArray[visitedTown]
			if j == 0 || numeratorDiffArray[visitedTown] < minimumNumeratorDiff {
				minimumNumeratorDiff = numeratorDiffArray[visitedTown]
			}
		}
	}

	fmt.Println("Diff array", numeratorDiffArray)
	for i := range numeratorDiffArray {
		if !ts.TownSlice[i].IsRequired && a.Visited[i] {
			numeratorDiffArray[i] -= minimumNumeratorDiff
			diffDenom += numeratorDiffArray[i]
		}
	}

	for i := range replacementNumeratorArray {
		if i == 0 {
			(*a).Probabilities[i] = numeratorDiffArray[i] / diffDenom
		} else {
			(*a).Probabilities[i] = (*a).Probabilities[i-1] + numeratorDiffArray[i]/diffDenom
		}
	}
}

func (a *Ant) substituteTown(ts Towns, unvisitedTownIndex int) {
	(*a).getReplaceProbabilityList(ts, unvisitedTownIndex)
	randFloat := randSource.Float64()
	replaceTownIndex := 0
	for randFloat > (*a).Probabilities[replaceTownIndex] {
		replaceTownIndex++
	}

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
