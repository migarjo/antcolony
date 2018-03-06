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
		if town.NormalizedRating != 0 {
			t.Error("NormalizedRating in town:", i, "equals", town.NormalizedRating, "but should be 0 when RatingPreference equals 0.")
		}
	}

	config.RatingPreference = 1

	towns.normalizeTownRatings(config)
	for i, town := range towns.TownSlice {
		if town.NormalizedRating != 0 {
			t.Error("NormalizedRating in town:", i, "equals", town.NormalizedRating, "but should be 0 when all Ratings submitted are equal.")
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
		if town.NormalizedRating != expectedNormalizedRating[i] {
			t.Error("NormalizedRating in town:", i, "equals", town.NormalizedRating, "but expected", expectedNormalizedRating[i])
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
