package io

import (
	"app/simulation"
	"app/simulationGraph"
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func ReadGraphFromFile(filepath string) simulationGraph.JsonGraphStructure {
	file, err := ioutil.ReadFile(filepath)
	if err != nil {
		fmt.Println("graph file reading problem")
	}

	gs := simulationGraph.JsonGraphStructure{}

	err = json.Unmarshal(file, &gs)

	if err != nil {
		fmt.Println("file content parsing to json problem")
		panic(err)
	}

	return gs
}

func SaveStatistics(filepath string, stats simulation.JsonStatsStructure) {
	jsonStats, err := json.Marshal(stats)

	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(filepath, jsonStats, 0644)
	if err != nil {
		fmt.Println("Could not write result to file.")
		return
	}
}

func SaveGraph(filepath string, g *simulationGraph.GraphWrapper) {
	jsonGraph, err := json.Marshal(simulationGraph.NewJsonGraphStructure(g))

	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(filepath, jsonGraph, 0644)
	if err != nil {
		fmt.Println("Could not save graph.")
		return
	}
}
