package main

// func TestCreateAntSlice(t *testing.T) {
// 	antRatio := 0.8
// 	inputJSON := readFixture("pentagon_input.json")
// 	config, towns, err := importInputs(inputJSON)
// 	if err != nil {
// 		check(err)
// 	}

// 	err = towns.initializeTowns(config)
// 	if err != nil {
// 		check(err)
// 	}

// 	numberOfAnts := int(math.Ceil(antRatio * float64(len(towns.TownSlice))))

// 	if numberOfAnts != 4 {
// 		t.Error("Expected numberOfAnts to be 4, got", numberOfAnts)
// 	}

// 	towns.calculateProbabilityMatrix(config)
// 	ants := createAntSlice(numberOfAnts, towns, config)

// 	if len(ants) != 4 {
// 		t.Error("Expected ant quantity to be 4, got", len(ants))
// 	}

// 	if len(ants[0].Tour) != 5 {
// 		t.Error("Expected number of towns visited to be 5, got", len(ants[0].Tour))
// 	}
// }
