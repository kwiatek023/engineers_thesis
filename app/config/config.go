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

	// Probability - specifies expression to evaluate probability of broken edge
	Probability string
}

// parseArgs - parses arguments passed in command line
func (args *AppArgs) parseArgs() {
	flag.StringVar(&args.GraphFile, "graph-file", "", "read graph structure from given file")
	flag.StringVar(&args.GraphType, "graph-type", "", "provide graph-type "+
		"(path,$number_of_vertices|clique,$number_of_vertices|regular,$number_of_vertices,$degree|grid,$height,$width|hypercube,$dimension"+
		"|tree,$number_of_vertices,$degree)")
	flag.BoolVar(&args.UseReliability, "use-reliability", false, "specifies if reliability of a network should be tested")
	flag.StringVar(&args.Probability, "p", "0.0", "specifies probability expression for reliability model")
}

// InitializeAppArgs - initializes and validates arguments
func InitializeAppArgs() AppArgs {
	var args = AppArgs{}
	args.parseArgs()
	flag.Parse()

	if args.GraphFile != "" && args.GraphType != "" {
		fmt.Println("You cannot use graph file while trying to build")
	} else if args.GraphFile == "" && args.GraphType == "" {
		fmt.Println("You have to specify graph file or graph type")
	}

	return args
}
