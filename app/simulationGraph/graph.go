package simulationGraph

import (
	"app/config"
	"github.com/yourbasic/graph"
)

type GraphWrapper struct {
	GraphStructure *graph.Mutable
	ReliabilityMap map[int]map[int]float64
}

func NewGraphWrapper(graphStructure *graph.Mutable, reliability map[int]map[int]float64) *GraphWrapper {
	return &GraphWrapper{GraphStructure: graphStructure, ReliabilityMap: reliability}
}

func BuildGraphFromConfig(conf config.JsonGraphStructure) *GraphWrapper {
	graphStructure := conf.Graph
	g := graph.New(graphStructure.NofVertices)

	var relMap = map[int]map[int]float64{}

	for _, e := range graphStructure.Edges {
		g.AddBoth(e.Edge[0], e.Edge[1])
		addReliability(relMap, e.Edge[0], e.Edge[1], e.Reliability)
	}

	return NewGraphWrapper(g, relMap)
}

func addReliability(relMap map[int]map[int]float64, firstVertex int, secondVertex int,
	rel float64) {
	relMap[firstVertex] = map[int]float64{}
	relMap[firstVertex][secondVertex] = rel
	relMap[secondVertex] = map[int]float64{}
	relMap[secondVertex][firstVertex] = rel
}
