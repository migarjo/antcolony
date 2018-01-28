package main

import (
	"math"
	"testing"
)

func TestInitializeTowns(t *testing.T) {
	//inputJSON := readFixture("input.json")
	//config, towns, err := importInputs(inputJSON)

}

func TestZeroScorePreference(t *testing.T) {
	inputJSON := readFixture("input.json")
	config, towns, err := importInputs(inputJSON)

	if err != nil {
		t.Error("Received error marshalling input: ", err)
	}

	towns.normalizeTownScores(config)
	for i, town := range towns.TownSlice {
		if town.NormalizedScore != 0 {
			t.Error("NormalizedScore in town:", i, "equals", town.NormalizedScore, "but should be 0 when ScorePreferene equals 0.")
		}
	}

	config.ScorePreference = 1

	towns.normalizeTownScores(config)
	for i, town := range towns.TownSlice {
		if town.NormalizedScore != 0 {
			t.Error("NormalizedScore in town:", i, "equals", town.NormalizedScore, "but should be 0 when all scores are equal.")
		}
	}

}

func TestNormalizeScores(t *testing.T) {
	inputJSON := readFixture("score_input.json")
	config, towns, err := importInputs(inputJSON)

	if err != nil {
		t.Error("Received error marshalling input: ", err)
	}

	towns.initializeTowns(config)

	expectedNormalizedScore := []float64{0, 0.5, 1}

	for i, town := range towns.TownSlice {
		if town.NormalizedScore != expectedNormalizedScore[i] {
			t.Error("NormalizedScore in town:", i, "equals", town.NormalizedScore, "but expected", expectedNormalizedScore[i])
		}
	}
}

func TestProbabilityMatrix(t *testing.T) {
	inputJSON := readFixture("score_input.json")
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
