package main

import (
	"io/ioutil"
	"testing"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}
func readFixture(fileName string) []byte {
	fixture, err := ioutil.ReadFile("test/fixtures/" + fileName)
	check(err)
	return fixture
}

func TestImportInputs(t *testing.T) {
	inputJSON := readFixture("input.json")
	config, towns, err := importInputs(inputJSON)

	if config.NumberOfIterations != 10 {
		t.Error("Expected Number of Iterations to be 10, got ", config.NumberOfIterations)
	}

	if len(towns.TownSlice) != 4 {
		t.Error("Expected 4 towns, got ", len(towns.TownSlice))
	}

	if !towns.TownSlice[0].IsRequired {
		t.Error("Expected Town 0 to have IsRequired be true because IncludesHome = true, got ", towns.TownSlice[0].IsRequired)
	}

	if towns.TownSlice[1].IsRequired {
		t.Error("Expected Town 1 to have IsRequired be false, got ", towns.TownSlice[1].IsRequired)
	}

	if !towns.TownSlice[3].IsRequired {
		t.Error("Expected Town 3 to have IsRequired be true, got ", towns.TownSlice[3].IsRequired)
	}

	if err != nil {
		t.Error("Received error marshalling input: ", err)
	}
	if config.VisitQuantity != 4 {
		t.Error("Expected visit quantity to be 4, got", config.VisitQuantity)
	}

	inputJSON = readFixture("pentagon_input.json")
	config, towns, err = importInputs(inputJSON)

	if config.VisitQuantity != 4 {
		t.Error("Expected visit quantity to be 4, got", config.VisitQuantity)
	}

}
