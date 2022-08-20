package main

import (
	"app/config"
	"app/simulation"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	manager := simulation.NewManager(config.InitializeAppArgs())
	manager.RunSimulation()

	//fmt.Println(g.GetDiameter())

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
