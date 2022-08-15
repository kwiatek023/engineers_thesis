package main

import (
	"app/config"
	"app/io"
	"app/simulation"
	"app/simulationGraph"
	"fmt"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	config.InitializeAppArgs()
	var g *simulationGraph.GraphWrapper
	if config.Args.GraphFile != "" {
		conf := io.ReadGraphFromFile(config.Args.GraphFile)
		g = simulationGraph.BuildGraphFromConfig(conf)
	} else {
		g = simulationGraph.BuildGraphFromType()
	}
	manager := simulation.NewManager(int(config.Args.NofVertices), g)
	manager.RunSimulation()

	fmt.Println(g)

	//g2 := simulationGraph.BuildGrid(2, 3, true)
	//fmt.Println(g2)
	//fmt.Println(g2.ReliabilityMap[0][3])
	//fmt.Println(g2.ReliabilityMap[3][0])
	//fmt.Println(simulationGraph.BuildDAryTree(4, 3, true))
	//g3 := simulationGraph.BuildCompleteGraph(4, true)
	//fmt.Println(g3)
	//fmt.Println(g3.ReliabilityMap[0][3])
	//fmt.Println(g3.ReliabilityMap[3][0])
	//fmt.Println(simulationGraph.BuildCompleteGraph(4, true))
	//fmt.Println(simulationGraph.BuildHyperCube(3, true))
	//fmt.Println(simulationGraph.BuildDRegularGraph(6, 4, true))
}
