package main

import (
	"app/config"
	"app/experiments"
	"app/io"
	"app/simulation"
	"app/simulationGraph"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	args := config.InitializeAppArgs()
	if args.Experiment == "" {
		fmt.Println("Building graph.")
		var g *simulationGraph.GraphWrapper
		if args.GraphFile != "" {
			conf := io.ReadGraphFromFile(args.GraphFile)
			g = simulationGraph.BuildGraphFromConfig(conf)
		} else {
			g = simulationGraph.BuildGraphFromType(args, true)
		}

		fmt.Println("Graph built.")
		fmt.Println("Simulation pending.")
		manager := simulation.NewManager(args.UseReliability, g)
		result := manager.RunSimulation(args.ProtocolName)

		//TODO add filenames options
		io.SaveStatistics("data.json", result)
		io.SaveGraph("graph.json", g)
		fmt.Println("Simulation finished.")
	} else {
		experiment := strings.Split(args.Experiment, ",")[0]
		if experiment == "extremaPropagation" {
			experiments.ExtremaPropagationExperiment(args.Experiment)
		} else if experiment == "countDistinct" {
			experiments.CountDistinctExperiment(args.Experiment)
		}
	}
}
