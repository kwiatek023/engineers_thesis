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

		if args.GraphCopyFile != "" {
			io.SaveGraph(args.GraphCopyFile, g)
		}
		fmt.Println("Simulation pending.")
		manager := simulation.NewManager(args.ReliabilityModel, g)
		result := manager.RunSimulation(args.ProtocolName)

		if args.StatsFile != "" {
			io.SaveStatistics(args.StatsFile, result)
		}
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
