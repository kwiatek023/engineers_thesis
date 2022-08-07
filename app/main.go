package main

import (
	"app/io"
	"app/simulation"
	"app/simulationGraph"
	"fmt"
)

func main() {
	fmt.Println("hello")

	conf := io.ReadGraphFromFile("./example.json")
	g := simulationGraph.BuildGraphFromConfig(conf)
	manager := simulation.NewManager(conf.Graph.NofVertices, 5, g)
	manager.RunSimulation()

}
