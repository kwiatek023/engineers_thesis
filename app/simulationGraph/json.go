package simulationGraph

// JsonGraphStructure - structure for storing graph data from file
type JsonGraphStructure struct {
	Graph JsonGraph `json:"graph"`
}

type JsonGraph struct {
	NofVertices uint       `json:"nofVertices"`
	Edges       []JsonEdge `json:"edges"`
}

type JsonEdge struct {
	Edge        []uint  `json:"edge"`
	Reliability float64 `json:"reliability"`
}

func NewJsonGraphStructure(g *GraphWrapper) *JsonGraphStructure {
	nofVertices := g.GraphStructure.Order()
	edges := g.GetEdges()
	jsonEdges := make([]JsonEdge, 0)
	for v, e := range edges {
		for w, _ := range e {
			jsonEdges = append(jsonEdges, JsonEdge{Edge: []uint{uint(v), uint(w)}, Reliability: g.GetRelMap()[v][w]})
		}
	}

	return &JsonGraphStructure{Graph: JsonGraph{NofVertices: uint(nofVertices), Edges: jsonEdges}}
}
