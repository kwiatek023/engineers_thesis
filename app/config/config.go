package config

import (
	"flag"
	"fmt"
)

// JsonGraphStructure - structure for storing graph data from file
type JsonGraphStructure struct {
	Graph struct {
		NofVertices uint `json:"nofVertices"`
		Edges       []struct {
			Edge        []uint  `json:"edge"`
			Reliability float64 `json:"reliability"`
		} `json:"edges"`
	} `json:"graph"`
}

// AppArgs - configuration of simulator
type AppArgs struct {
	// GraphFile - if provided, program uses graph topology from file
	GraphFile string

	// GraphType - specifies graph topology
	GraphType string

	// NofVertices - number of vertices
	NofVertices uint

	// UseReliability - specifies if reliability of a network should be tested
	UseReliability bool
}

// parseArgs - parses arguments passed in command line
func (args *AppArgs) parseArgs() {
	flag.StringVar(&args.GraphFile, "graph-file", "", "read graph structure from given file")
	flag.UintVar(&args.NofVertices, "n", 0, "provide number of stations")
	flag.StringVar(&args.GraphType, "graph-type", "", "provide graph-type "+
		"(path|clique|regular,$degree|grid,$height,$width|hypercube|tree,$degree)")
	flag.BoolVar(&args.UseReliability, "use-reliability", false, "")
}

// InitializeAppArgs - initializes and validates arguments
func InitializeAppArgs() AppArgs {
	var args = AppArgs{}
	args.parseArgs()

	if args.GraphFile != "" && args.GraphType != "" {
		fmt.Println("You cannot use graph file while trying to build", args.GraphType)
	} else if args.GraphFile == "" && args.GraphType == "" {
		fmt.Println("You have to specify graph file or graph type")
	}

	if args.NofVertices == 0 {
		fmt.Println("Number of vertices must be greater than 0")
	}

	flag.Parse()

	return args
}
