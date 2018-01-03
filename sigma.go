package antcolony

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

type SigmaObject struct {
	Nodes []SigmaNode `json:"nodes"`
	Edges []SigmaEdge `json:"edges"`
}

type SigmaNode struct {
	ID    string `json:"id"`
	Label string `json:"label"`
	X     string `json:"x"`
	Y     string `json:"y"`
	Size  string `json:"size"`
}

type SigmaEdge struct {
	ID     string `json:"id"`
	Source string `json:"source"`
	Target string `json:"target"`
	Color  string `json:"color,omitempty"`
	Size   string `json:"size,omitempty"`
}

func createSigmaObject(ts *Towns, a *ant) SigmaObject {
	so := SigmaObject{
		Nodes: []SigmaNode{},
		Edges: []SigmaEdge{},
	}

	if ts != nil {
		so.addSigmaNodes(*ts)
	}

	if a != nil {
		so.addSigmaEdges(*a)
	}
	return so
}

func (o *SigmaObject) addSigmaNodes(ts Towns) {
	for _, t := range ts.TownSlice {
		node := SigmaNode{
			ID:    strconv.Itoa(t.ID),
			Label: strconv.Itoa(t.ID) + " (" + strconv.Itoa(t.XCoord) + "," + strconv.Itoa(t.YCoord) + ")",
			X:     strconv.Itoa(t.XCoord),
			Y:     strconv.Itoa(t.YCoord),
			Size:  "1",
		}
		(*o).Nodes = append((*o).Nodes, node)
	}
}

func (o *SigmaObject) addSigmaEdges(a ant) {
	fmt.Println("BestAnt:")
	a.printAnt()
	for i := 1; i < len(a.tour); i++ {

		edge := SigmaEdge{
			ID:     strconv.Itoa(i),
			Source: strconv.Itoa(a.tour[i-1]),
			Target: strconv.Itoa(a.tour[i]),
			Size:   "1",
		}
		(*o).Edges = append((*o).Edges, edge)
	}
}

func (o *SigmaObject) jsonify() []byte {
	soJSON, err := json.Marshal(*o)

	if err != nil {
		fmt.Println(err)
	}
	return soJSON
}

func (o *SigmaObject) writeToFile(path string) {
	soJSON := (*o).jsonify()

	f, err := os.Create(path)
	if err != nil {
		fmt.Println(err)
	}

	n, err := f.Write(soJSON)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Wrote:", n, "objects")
}
