package main

import (
	"math"
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
		if town.NormalizedRating != 1 {
			t.Error("NormalizedRating in town:", i, "equals", town.NormalizedRating, "but should be 1 when RatingPreference equals 0.")
		}
	}

	config.RatingPreference = 1

	towns.normalizeTownRatings(config)
	for i, town := range towns.TownSlice {
		if town.NormalizedRating != 1 {
			t.Error("NormalizedRating in town:", i, "equals", town.NormalizedRating, "but should be 1 when all Ratings submitted are equal.")
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

	expectedProbabilityMatrix := [][]float64{[]float64{math.Inf(1), 1, 0.5}, []float64{1.5, math.Inf(1), 1.5}, []float64{1.5, 2, math.Inf(1)}}

	for i := range towns.ProbabilityMatrix {
		for j := range towns.ProbabilityMatrix[i] {
			if towns.ProbabilityMatrix[i][j] != expectedProbabilityMatrix[i][j] {
				t.Error("ProbabilityMatrix for i:", i, "j:", j, "equals", towns.ProbabilityMatrix[i][j], "but expected", expectedProbabilityMatrix[i][j])
			}
		}

	}
}
