package simulationGraph

import (
	"app/config"
	"github.com/yourbasic/graph"
	"github.com/yourbasic/graph/build"
	"math"
	"math/rand"
)

type GraphWrapper struct {
	GraphStructure *graph.Mutable
	ReliabilityMap map[int]map[int]float64
}

type nothing struct{}

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

func BuildGrid(m, n int, useReliability bool) *GraphWrapper {
	nofVertices := m * n
	g := graph.New(nofVertices)
	virtualGrid := build.Grid(m, n)

	// set of edges
	var edges = map[int]map[int]nothing{}
	for i := 0; i < nofVertices; i++ {
		edges[i] = map[int]nothing{}
	}

	for i := 0; i < nofVertices; i++ {
		virtualGrid.Visit(i, func(w int, c int64) (skip bool) {
			g.AddBoth(i, w)
			_, e1 := edges[i][w]
			_, e2 := edges[w][i]

			//fmt.Println(e1, e2)
			if !(e1 || e2) {
				edges[i][w] = nothing{}
			}

			return
		})
	}

	if !useReliability {
		return NewGraphWrapper(g, nil)
	}

	var relMap = map[int]map[int]float64{}
	maxRel := 1.0
	minRel := 0.9

	for v, e := range edges {
		for w, _ := range e {
			rel := minRel + rand.Float64()*(maxRel-minRel)
			addReliability(relMap, v, w, rel)
		}
	}

	return NewGraphWrapper(g, relMap)
}

func BuildDAryTree(nofVertices, degree int, useReliability bool) *GraphWrapper {
	levels := int(math.Ceil(math.Log(float64(nofVertices*(degree-1)+1)) / math.Log(float64(degree))))
	virtualGrid := build.Tree(degree, levels)
	g := graph.New(nofVertices)

	// set of edges
	var edges = map[int]map[int]nothing{}
	for i := 0; i < nofVertices; i++ {
		edges[i] = map[int]nothing{}
	}

	for i := 0; i < nofVertices; i++ {
		virtualGrid.Visit(i, func(w int, c int64) (skip bool) {
			if w < nofVertices {
				g.AddBoth(i, w)
				_, e1 := edges[i][w]
				_, e2 := edges[w][i]

				//fmt.Println(e1, e2)
				if !(e1 || e2) {
					edges[i][w] = nothing{}
				}
			}

			return
		})
	}

	if !useReliability {
		return NewGraphWrapper(g, nil)
	}

	var relMap = map[int]map[int]float64{}
	maxRel := 1.0
	minRel := 0.9

	for v, e := range edges {
		for w, _ := range e {
			rel := minRel + rand.Float64()*(maxRel-minRel)
			addReliability(relMap, v, w, rel)
		}
	}

	return NewGraphWrapper(g, relMap)
}
