package simulationGraph

import (
	"app/config"
	"fmt"
	"github.com/yourbasic/graph"
	"github.com/yourbasic/graph/build"
	"math"
	"math/rand"
	"strconv"
	"strings"
)

type GraphWrapper struct {
	GraphStructure *graph.Mutable
	ReliabilityMap map[int]map[int]float64
	diameter       int
}

type nothing struct{}

func NewGraphWrapper(graphStructure *graph.Mutable, reliability map[int]map[int]float64) *GraphWrapper {
	return &GraphWrapper{GraphStructure: graphStructure, ReliabilityMap: reliability, diameter: 0}
}

func BuildGraphFromConfig(conf config.JsonGraphStructure) *GraphWrapper {
	graphStructure := conf.Graph
	g := graph.New(int(graphStructure.NofVertices))

	var relMap = initRelMap(int(graphStructure.NofVertices))

	for _, e := range graphStructure.Edges {
		g.AddBoth(int(e.Edge[0]), int(e.Edge[1]))
		addReliability(relMap, int(e.Edge[0]), int(e.Edge[1]), e.Reliability)
	}

	return NewGraphWrapper(g, relMap)
}

func initRelMap(nofVertices int) map[int]map[int]float64 {
	var relMap = map[int]map[int]float64{}
	for i := 0; i < nofVertices; i++ {
		relMap[i] = map[int]float64{}
	}

	return relMap
}

func addReliability(relMap map[int]map[int]float64, firstVertex int, secondVertex int,
	rel float64) {
	relMap[firstVertex][secondVertex] = rel
	relMap[secondVertex][firstVertex] = rel
}

func BuildPath(nofVertices int, useReliability bool) *GraphWrapper {
	g := graph.New(nofVertices)

	for i := 0; i < nofVertices-1; i++ {
		g.AddBoth(i, i+1)
	}

	if !useReliability {
		return NewGraphWrapper(g, nil)
	}

	maxRel := 1.0
	minRel := 0.9
	var relMap = map[int]map[int]float64{}

	for i := 0; i < nofVertices-1; i++ {
		rel := minRel + rand.Float64()*(maxRel-minRel)
		addReliability(relMap, i, i+1, rel)
	}

	return NewGraphWrapper(g, relMap)
}

func BuildCompleteGraph(nofVertices int, useReliability bool) *GraphWrapper {
	virtualCompleteGraph := build.Kn(nofVertices)
	return convertVirtualToMutable(nofVertices, useReliability, virtualCompleteGraph)
}

func BuildGrid(m, n int, useReliability bool) *GraphWrapper {
	nofVertices := m * n
	virtualGrid := build.Grid(m, n)

	return convertVirtualToMutable(nofVertices, useReliability, virtualGrid)
}

func BuildDAryTree(nofVertices, degree int, useReliability bool) *GraphWrapper {
	levels := int(math.Ceil(math.Log(float64(nofVertices*(degree-1)+1)) / math.Log(float64(degree))))
	virtualTree := build.Tree(degree, levels)

	return convertVirtualToMutable(nofVertices, useReliability, virtualTree)
}

func BuildDRegularGraph(nofVertices, degree int, useReliability bool) *GraphWrapper {
	if !(degree < nofVertices && (nofVertices%2 == 0 || degree%2 == 0)) {
		fmt.Println("Improper data for d-regular graph")
		//	raise error
		return nil
	}

	m := degree / 2
	distancesBetweenVertexNeighbours := make([]int, 0)
	for i := 1; i <= m; i++ {
		distancesBetweenVertexNeighbours = append(distancesBetweenVertexNeighbours, i)
	}

	if nofVertices%2 == 0 && degree%2 != 0 {
		distancesBetweenVertexNeighbours = append(distancesBetweenVertexNeighbours, nofVertices/2)
	}
	dRegularVirtualGraph := build.Circulant(nofVertices, distancesBetweenVertexNeighbours...)

	return convertVirtualToMutable(nofVertices, useReliability, dRegularVirtualGraph)
}

func BuildHyperCube(dimensions int, useReliability bool) *GraphWrapper {
	virtualHyperCube := build.Hyper(dimensions)
	nofVertices := int(math.Pow(2, float64(dimensions)))

	return convertVirtualToMutable(nofVertices, useReliability, virtualHyperCube)
}

func convertVirtualToMutable(nofVertices int, useReliability bool, immutableGraph *build.Virtual) *GraphWrapper {
	g := graph.New(nofVertices)

	// set of edges
	var edges = map[int]map[int]nothing{}
	for i := 0; i < nofVertices; i++ {
		edges[i] = map[int]nothing{}
	}

	for i := 0; i < nofVertices; i++ {
		immutableGraph.Visit(i, func(w int, c int64) (skip bool) {
			if w < nofVertices {
				g.AddBoth(i, w)
				_, e1 := edges[i][w]
				_, e2 := edges[w][i]

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

	var relMap = initRelMap(nofVertices)
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

func BuildGraphFromType(args config.AppArgs) *GraphWrapper {
	params := strings.Split(args.GraphType, ",")
	graphName := params[0]
	var g *GraphWrapper

	switch strings.ToLower(graphName) {
	case "clique":
		{
			g = BuildCompleteGraph(int(args.NofVertices), args.UseReliability)
		}
	case "hypercube":
		{
			dimensions, err := strconv.Atoi(params[1])

			if err != nil {
				// TODO handle error
				panic(err)
			}

			g = BuildHyperCube(dimensions, args.UseReliability)
		}
	case "path":
		{
			g = BuildPath(int(args.NofVertices), args.UseReliability)
		}
	case "grid":
		{
			m, err := strconv.Atoi(params[1])
			if err != nil {
				// TODO handle error
				panic(err)
			}

			n, err := strconv.Atoi(params[2])
			if err != nil {
				// TODO handle error
				panic(err)
			}

			g = BuildGrid(m, n, args.UseReliability)
		}
	case "tree":
		{
			degree, err := strconv.Atoi(params[1])
			if err != nil {
				// TODO handle error
				panic(err)
			}

			g = BuildDAryTree(int(args.NofVertices), degree, args.UseReliability)
		}
	case "regular":
		{
			degree, err := strconv.Atoi(params[1])
			if err != nil {
				// TODO handle error
				panic(err)
			}

			g = BuildDRegularGraph(int(args.NofVertices), degree, args.UseReliability)
		}
	}

	return g
}

func allPairsShortestPath(g *GraphWrapper) map[int]map[int]float64 {
	var dist = map[int]map[int]float64{}
	nofVertices := g.GraphStructure.Order()

	for i := 0; i < nofVertices; i++ {
		dist[i] = map[int]float64{}
	}

	for i := 0; i < nofVertices; i++ {
		dist[i][i] = 0
	}

	for i := 0; i < nofVertices; i++ {
		for j := 0; j < nofVertices; j++ {
			if i != j {
				if g.GraphStructure.Edge(i, j) {
					dist[i][j] = 1
				} else {
					dist[i][j] = math.Inf(1)
				}
			}
		}
	}

	for k := 0; k < nofVertices; k++ {
		for i := 0; i < nofVertices; i++ {
			for j := 0; j < nofVertices; j++ {
				if dist[i][k]+dist[k][j] < dist[i][j] {
					dist[i][j] = dist[i][k] + dist[k][j]
				}
			}
		}
	}

	return dist
}

func calcDiameter(g *GraphWrapper) float64 {
	nofVertices := g.GraphStructure.Order()
	dist := allPairsShortestPath(g)
	max := 0.0

	for i := 0; i < nofVertices; i++ {
		for j := 0; j < nofVertices; j++ {
			if dist[i][j] != math.Inf(1) {
				if dist[i][j] > max {
					max = dist[i][j]
				}
			}
		}
	}

	return max
}

func (g *GraphWrapper) GetDiameter() int {
	if g.diameter == 0 {
		g.diameter = int(calcDiameter(g))
	}

	return g.diameter
}
