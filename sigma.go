package main

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// SigmaEdge is the struct that represents each edge between towns in sigma.js syntax
type SigmaEdge struct {
	ID     string `json:"id"`
	Source string `json:"source"`
	Target string `json:"target"`
	Color  string `json:"color,omitempty"`
	Size   string `json:"size,omitempty"`
}

func exportSigmaEdges(a ant) string {

	edges := []SigmaEdge{}
	for i := 1; i < len(a.tour); i++ {

		edge := SigmaEdge{
			ID:     strconv.Itoa(i),
			Source: strconv.Itoa(a.tour[i-1]),
			Target: strconv.Itoa(a.tour[i]),
			Size:   "1",
		}
		edges = append(edges, edge)
	}

	edgeJSON, err := json.Marshal(edges)

	if err != nil {
		fmt.Println(err)
	}
	return string(edgeJSON[:])
}
