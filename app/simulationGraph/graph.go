package simulationGraph

import (
	"app/config"
	"app/utils"
	"github.com/yourbasic/graph"
	"github.com/yourbasic/graph/build"
	"log"
	"math"
	"strings"
)

// GraphWrapper - graph structure used in simulations
type GraphWrapper struct {
	GraphStructure *graph.Mutable
	reliabilityMap map[int]map[int]float64
	diameter       int
	edges          map[int]map[int]nothing
}

type nothing struct{}

func NewGraphWrapper(graphStructure *graph.Mutable, reliability map[int]map[int]float64, edges map[int]map[int]nothing) *GraphWrapper {
	return &GraphWrapper{GraphStructure: graphStructure, reliabilityMap: reliability, diameter: 0, edges: edges}
}

// BuildGraphFromConfig - function which builds graph based on given file
func BuildGraphFromConfig(conf JsonGraphStructure) *GraphWrapper {
	graphStructure := conf.Graph
	g := graph.New(int(graphStructure.NofVertices))

	var relMap = initRelMap(int(graphStructure.NofVertices))

	for _, e := range graphStructure.Edges {
		g.AddBoth(int(e.Edge[0]), int(e.Edge[1]))
		addReliability(relMap, int(e.Edge[0]), int(e.Edge[1]), e.Reliability)
	}

	edges := makeEdgeSet(g)
	resultGraph := NewGraphWrapper(g, relMap, edges)
	resultGraph.diameter = int(calcDiameter(resultGraph))
	return resultGraph
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

func BuildPath(nofVertices int, reliabilityModel string, pExpr string) *GraphWrapper {
	g := graph.New(nofVertices)

	for i := 0; i < nofVertices-1; i++ {
		g.AddBoth(i, i+1)
	}

	edges := makeEdgeSet(g)

	if reliabilityModel == "" {
		return NewGraphWrapper(g, nil, edges)
	}

	parameters := make(map[string]interface{}, 0)
	parameters["n"] = nofVertices
	p := utils.EvaluateExpression(pExpr, parameters)
	var relMap = initRelMap(g.Order())

	for i := 0; i < nofVertices-1; i++ {
		addReliability(relMap, i, i+1, p)
	}

	return NewGraphWrapper(g, relMap, edges)
}

func BuildClique(nofVertices int, reliabilityModel string, pExpr string) *GraphWrapper {
	parameters := make(map[string]interface{}, 0)
	parameters["n"] = nofVertices
	p := utils.EvaluateExpression(pExpr, parameters)

	virtualCompleteGraph := build.Kn(nofVertices)
	return convertVirtualToMutable(nofVertices, reliabilityModel, virtualCompleteGraph, p)
}

func BuildGrid(m, n int, reliabilityModel string, pExpr string) *GraphWrapper {
	nofVertices := m * n
	virtualGrid := build.Grid(m, n)
	parameters := make(map[string]interface{}, 0)
	parameters["n"] = nofVertices
	p := utils.EvaluateExpression(pExpr, parameters)

	return convertVirtualToMutable(nofVertices, reliabilityModel, virtualGrid, p)
}

func BuildDAryTree(nofVertices, degree int, reliabilityModel string, pExpr string) *GraphWrapper {
	levels := int(math.Ceil(math.Log(float64(nofVertices*(degree-1)+1)) / math.Log(float64(degree))))
	virtualTree := build.Tree(degree, levels)
	parameters := make(map[string]interface{}, 0)
	parameters["n"] = nofVertices
	p := utils.EvaluateExpression(pExpr, parameters)

	return convertVirtualToMutable(nofVertices, reliabilityModel, virtualTree, p)
}

func BuildDRegularGraph(nofVertices, degree int, reliabilityModel string, pExpr string) *GraphWrapper {
	if !(degree < nofVertices && (nofVertices%2 == 0 || degree%2 == 0)) {
		log.Fatal("Improper data for d-regular graph")
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
	parameters := make(map[string]interface{}, 0)
	parameters["n"] = nofVertices
	p := utils.EvaluateExpression(pExpr, parameters)

	return convertVirtualToMutable(nofVertices, reliabilityModel, dRegularVirtualGraph, p)
}

func BuildHyperCube(dimensions int, reliabilityModel string, pExpr string) *GraphWrapper {
	virtualHyperCube := build.Hyper(dimensions)
	nofVertices := int(math.Pow(2, float64(dimensions)))
	parameters := make(map[string]interface{}, 0)
	parameters["n"] = nofVertices
	p := utils.EvaluateExpression(pExpr, parameters)

	return convertVirtualToMutable(nofVertices, reliabilityModel, virtualHyperCube, p)
}

func BuildGridOfCliques(m, n, nofVerticesInClique int, reliabilityModel string, pExpr string) *GraphWrapper {
	nofVertices := m * n * nofVerticesInClique
	clique := build.Kn(nofVerticesInClique)
	row := build.Kn(nofVerticesInClique)

	parameters := make(map[string]interface{}, 0)
	parameters["n"] = nofVertices
	p := utils.EvaluateExpression(pExpr, parameters)

	// create row
	for i := 0; i < n-1; i++ {
		from := build.Vertex(i * nofVerticesInClique)
		to := build.Vertex((i + 1) * nofVerticesInClique)
		row = row.Join(clique, build.EdgeSet{From: from, To: to})
	}

	rowSize := n * nofVerticesInClique
	result := build.Empty(1)

	result = result.Connect(0, row)
	// create columns
	for i := 0; i < m-1; i++ {
		from := build.Vertex(i * rowSize)
		to := build.Vertex((i + 1) * rowSize)
		result = result.Join(row, build.EdgeSet{From: from, To: to})
	}

	//	add missing edges
	for j := 0; j < m-1; j++ {
		for i := 1; i < n; i++ {
			from := build.Vertex(j*rowSize + i*nofVerticesInClique)
			to := build.Vertex(j*rowSize + i*nofVerticesInClique + rowSize)
			result = result.Add(build.EdgeSet{From: from, To: to})
		}
	}

	res := convertVirtualToMutable(nofVertices, reliabilityModel, result, p)
	return res
}

func makeEdgeSet(g *graph.Mutable) map[int]map[int]nothing {
	nofVertices := g.Order()
	var edges = map[int]map[int]nothing{}
	for i := 0; i < nofVertices; i++ {
		edges[i] = map[int]nothing{}
	}

	for i := 0; i < nofVertices; i++ {
		g.Visit(i, func(w int, c int64) (skip bool) {
			if w < nofVertices {
				_, e1 := edges[i][w]
				_, e2 := edges[w][i]

				if !(e1 || e2) {
					edges[i][w] = nothing{}
				}
			}

			return
		})
	}

	return edges
}

func convertVirtualToMutable(nofVertices int, reliabilityModel string, virtualGraph *build.Virtual, p float64) *GraphWrapper {
	g := graph.New(nofVertices)

	for i := 0; i < nofVertices; i++ {
		virtualGraph.Visit(i, func(w int, c int64) (skip bool) {
			if w < nofVertices {
				g.AddBoth(i, w)
			}

			return
		})
	}

	edges := makeEdgeSet(g)

	if reliabilityModel == "" {
		return NewGraphWrapper(g, nil, edges)
	}

	var relMap = initRelMap(nofVertices)

	for v, e := range edges {
		for w, _ := range e {
			addReliability(relMap, v, w, p)
		}
	}

	return NewGraphWrapper(g, relMap, edges)
}

// BuildGraphFromType - function which builds graph based on given type
func BuildGraphFromType(args config.AppArgs, buildWithDiameter bool) *GraphWrapper {
	params := strings.Split(args.GraphType, ",")
	graphName := params[0]
	var g *GraphWrapper

	switch strings.ToLower(graphName) {
	case "clique":
		{
			nofVertices := utils.ParseStrToPositiveInt(params[1])
			g = BuildClique(nofVertices, args.ReliabilityModel, args.Probability)
		}
	case "hypercube":
		{
			dimensions := utils.ParseStrToPositiveInt(params[1])
			g = BuildHyperCube(dimensions, args.ReliabilityModel, args.Probability)
		}
	case "path":
		{
			nofVertices := utils.ParseStrToPositiveInt(params[1])
			g = BuildPath(nofVertices, args.ReliabilityModel, args.Probability)
		}
	case "grid":
		{
			m := utils.ParseStrToPositiveInt(params[1])
			n := utils.ParseStrToPositiveInt(params[2])
			g = BuildGrid(m, n, args.ReliabilityModel, args.Probability)
		}
	case "tree":
		{
			nofVertices := utils.ParseStrToPositiveInt(params[1])
			degree := utils.ParseStrToPositiveInt(params[2])
			g = BuildDAryTree(nofVertices, degree, args.ReliabilityModel, args.Probability)
		}
	case "regular":
		{
			nofVertices := utils.ParseStrToPositiveInt(params[1])
			degree := utils.ParseStrToPositiveInt(params[2])
			g = BuildDRegularGraph(nofVertices, degree, args.ReliabilityModel, args.Probability)
		}
	case "gridOfCliques":
		{
			m := utils.ParseStrToPositiveInt(params[1])
			n := utils.ParseStrToPositiveInt(params[2])
			nofVerticesInClique := utils.ParseStrToPositiveInt(params[3])
			g = BuildGridOfCliques(m, n, nofVerticesInClique, args.ReliabilityModel, args.Probability)
		}
	}

	if buildWithDiameter {
		g.diameter = int(calcDiameter(g))
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

func (g *GraphWrapper) SetDiameter(diameter int) {
	g.diameter = diameter
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
	return g.diameter
}

func (g *GraphWrapper) GetEdges() map[int]map[int]nothing {
	return g.edges
}

func (g *GraphWrapper) GetRelMap() map[int]map[int]float64 {
	return g.reliabilityMap
}
