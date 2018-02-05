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
