package antcolony

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
)

type ant struct {
	id            int
	tour          []int
	visited       []bool
	probabilities []float64
	score         float64
}

type ProgressOverTime struct {
	Generation   []int     `json:"labels"`
	AverageScore []float64 `json:"average"`
	MinimumScore []float64 `json:"minimum"`
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

func (a *ant) getProbabilityList(ts Towns) {
	// Current location
	i := (*a).tour[len((*a).tour)-1]
	t := ts.TownSlice[i]
	n := len(ts.TownSlice)

	denom := 0.0
	numerator := make([]float64, n)

	for l := 0; l < n; l++ {
		if !(*a).visited[l] {
			if ts.ProbabilityMatrix[i][l] != 0 {
				numerator[l] = ts.ProbabilityMatrix[i][l]
				denom += numerator[l]
			} else {
				numerator[l] = math.Pow(t.Trails[l], trailPreference) * math.Pow(1.0/t.Distances[l], distancePreference)
				ts.ProbabilityMatrix[i][l] = numerator[l]
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

func (a *ant) visitTown(t Town, ts []Town) {
	(*a).tour = append((*a).tour, t.ID)
	(*a).visited[t.ID] = true
	if len((*a).tour) > 1 {
		(*a).score += getDistance(t, ts[(*a).tour[len((*a).tour)-2]])
	}
}

func (a *ant) visitNextTown(ts Towns) {
	(*a).getProbabilityList(ts)
	randFloat := randSource.Float64()
	i := 0
	for randFloat > (*a).probabilities[i] {
		i++
	}
	(*a).tour = append((*a).tour, i)
	(*a).visited[i] = true
	(*a).score += getDistance(ts.TownSlice[i], ts.TownSlice[(*a).tour[len((*a).tour)-2]])
}

func (a ant) printAnt() {
	fmt.Println(a.tour, a.score)
}

func printAnts(a []ant) {
	for _, ant := range a {
		ant.printAnt()
	}
}

func createAntSlice(n int, ts Towns) []ant {
	ants := []ant{}

	for a := 0; a < n; a++ {
		myAnt := createAnt(a, len(ts.TownSlice))
		myAnt.visitTown(ts.TownSlice[0], ts.TownSlice)

		for len(myAnt.tour) < len(ts.TownSlice) {
			myAnt.visitNextTown(ts)
		}
		myAnt.visitTown(ts.TownSlice[0], ts.TownSlice)

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

func (p *ProgressOverTime) add(averageScore float64, minimumScore float64) {
	(*p).Generation = append((*p).Generation, len((*p).Generation))
	(*p).AverageScore = append((*p).AverageScore, averageScore)
	(*p).MinimumScore = append((*p).MinimumScore, minimumScore)
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
