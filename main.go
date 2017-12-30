package main

import (
	"fmt"
	"math/rand"
	"time"
)

type model struct{}

type node struct {
	id     int
	xCoord int
	yCoord int
}

func main() {
	home := createNodes(1, 100)
	nodes := createNodes(10, 100)

	fmt.Println(home)
	fmt.Println(nodes)
}

func createNodes(n int, fieldSize int) []node {
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)

	nodeSlice := []node{}
	for i := 0; i < n; i++ {
		r.Intn(fieldSize - 1)
		thisNode := node{i, r.Intn(fieldSize - 1), r.Intn(fieldSize - 1)}
		nodeSlice = append(nodeSlice, thisNode)
	}
	return nodeSlice
}
