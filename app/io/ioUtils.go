package io

import (
	"app/config"
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
	}

	return gs
}
