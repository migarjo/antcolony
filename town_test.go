package main

import (
	"testing"
)

func TestInitializeTowns(t *testing.T) {
	//inputJSON := readFixture("input.json")
	//config, towns, err := importInputs(inputJSON)

}

func TestZeroRatingPreference(t *testing.T) {
	inputJSON := readFixture("input.json")
	config, towns, err := importInputs(inputJSON)

	if err != nil {
		t.Error("Received error marshalling input: ", err)
	}

	towns.normalizeTownRatings(config)
	for i, town := range towns.TownSlice {
		if town.normalizedRating != 0 {
			t.Error("NormalizedRating in town:", i, "equals", town.normalizedRating, "but should be 0 when RatingPreference equals 0.")
		}
	}

	config.RatingPreference = 1

	towns.normalizeTownRatings(config)
	for i, town := range towns.TownSlice {
		if town.normalizedRating != 0 {
			t.Error("NormalizedRating in town:", i, "equals", town.normalizedRating, "but should be 0 when all Ratings submitted are equal.")
		}
	}

}

func TestNormalizeRatings(t *testing.T) {
	inputJSON := readFixture("rating_input.json")
	config, towns, err := importInputs(inputJSON)

	if err != nil {
		t.Error("Received error marshalling input: ", err)
	}

	towns.initializeTowns(config)

	expectedNormalizedRating := []float64{0, 0.5, 1}

	for i, town := range towns.TownSlice {
		if town.normalizedRating != expectedNormalizedRating[i] {
			t.Error("NormalizedRating in town:", i, "equals", town.normalizedRating, "but expected", expectedNormalizedRating[i])
		}
	}
}

func TestProbabilityMatrix(t *testing.T) {
	inputJSON := readFixture("rating_input.json")
	config, towns, err := importInputs(inputJSON)

	if err != nil {
		t.Error("Received error marshalling input: ", err)
	}

	towns.initializeTowns(config)
	towns.calculateProbabilityMatrix(config)

	expectedProbabilityMatrix := [][]float64{[]float64{0, 1, 0.5}, []float64{1.5, 0, 1.5}, []float64{1.5, 2, 0}}

	for i := range towns.probabilityMatrix {
		for j := range towns.probabilityMatrix[i] {
			if towns.probabilityMatrix[i][j] != expectedProbabilityMatrix[i][j] {
				t.Error("ProbabilityMatrix for i:", i, "j:", j, "equals", towns.probabilityMatrix[i][j], "but expected", expectedProbabilityMatrix[i][j])
			}
		}
	}
}

func TestIsAvailable(t *testing.T) {

	towns := Towns{
		TownSlice: []Town{
			Town{
				AvailabilityBounds: AvailabilityBounds{
					Start: 1,
					End:   4,
				},
				VisitDuration: 1,
				Distances: []float64{
					0,
					1.0,
					2.0,
					3.0,
					4.0,
				},
			},
			Town{
				AvailabilityBounds: AvailabilityBounds{
					Start: 0,
					End:   0,
				},
				VisitDuration: 1,
				Distances: []float64{
					1.0,
					0,
					1.0,
					2.0,
					3.0,
				},
			},
			Town{
				AvailabilityBounds: AvailabilityBounds{
					Start: 0,
					End:   2,
				},
				VisitDuration: 1,
				Distances: []float64{
					2.0,
					1.0,
					0,
					1.0,
					2.0,
				},
			},
			Town{
				AvailabilityBounds: AvailabilityBounds{
					Start: 4,
					End:   6,
				},
				VisitDuration: 1,
				Distances: []float64{
					3.0,
					2.0,
					1.0,
					0,
					1.0,
				},
			},
			Town{
				AvailabilityBounds: AvailabilityBounds{
					Start: 0,
					End:   0,
				},
				VisitDuration: 1,
				Distances: []float64{
					4.0,
					3.0,
					2.0,
					1.0,
					0,
				},
			},
		},
	}

	ant := Ant{
		Tour: []int{
			1,
		},
		VisitSpan: [][]float64{
			[]float64{0, 1},
		},
		costState: 1,
	}

	if isAvailable(towns, &ant, 2) {
		t.Error("Expected town not to be available when the cost of the visit is larger than the town's End")
	}

	if !isAvailable(towns, &ant, 0) {
		t.Error("Expected town to be available when the cost of the visit is between or equal to Start and End")
	}

	if isAvailable(towns, &ant, 3) {
		t.Error("Expected town not to be available when the cost of the visit is smaller than the town's Start")
	}

	if !isAvailable(towns, &ant, 4) {
		t.Error("Expected town to be available when Start and End are 0")
	}
}
