package io

import (
	"app/config"
	"app/simulation"
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func ReadGraphFromFile(filepath string) config.JsonGraphStructure {
	file, err := ioutil.ReadFile(filepath)
	if err != nil {
		fmt.Println("graph file reading problem")
	}

	gs := config.JsonGraphStructure{}

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
