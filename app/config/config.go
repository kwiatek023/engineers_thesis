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

	// UseReliability - specifies if reliability of a network should be tested
	UseReliability bool

	// Probability - specifies probability of broken edge
	Probability float64
}

// parseArgs - parses arguments passed in command line
func (args *AppArgs) parseArgs() {
	flag.StringVar(&args.GraphFile, "graph-file", "", "read graph structure from given file")
	flag.StringVar(&args.GraphType, "graph-type", "", "provide graph-type "+
		"(path,$number_of_vertices|clique,$number_of_vertices|regular,$number_of_vertices,$degree|grid,$height,$width|hypercube,$dimension"+
		"|tree,$number_of_vertices,$degree)")
	flag.BoolVar(&args.UseReliability, "use-reliability", false, "specifies if reliability of a network should be tested")
	flag.Float64Var(&args.Probability, "p", 0, "specifies probability of broken edge if graph type is provided")
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

	if args.Probability < 0 || args.Probability > 1 {
		fmt.Println("Probability should be in range 0 .. 1")
	}

	flag.Parse()

	return args
}
