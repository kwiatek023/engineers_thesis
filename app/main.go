package main

import (
	"app/config"
	"app/io"
	"app/simulation"
	"app/simulationGraph"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	args := config.InitializeAppArgs()
	var g *simulationGraph.GraphWrapper
	if args.GraphFile != "" {
		conf := io.ReadGraphFromFile(args.GraphFile)
		g = simulationGraph.BuildGraphFromConfig(conf)
	} else {
		g = simulationGraph.BuildGraphFromType(args)
	}

	manager := simulation.NewManager(args, g)
	result := manager.RunSimulation()

	io.SaveStatistics("data.json", result)
}
