package config

import (
	"flag"
	"fmt"
)

// AppArgs - configuration of simulator
// TODO add -help which informs how to use the app
type AppArgs struct {
	// GraphFile - if provided, program uses graph topology from a file
	GraphFile string

	// GraphCopyFile - if provided, program saves copy of a graph to a file
	GraphCopyFile string

	// StatsFile - if provided, program saves statistics to a file
	StatsFile string

	// GraphType - specifies graph topology
	GraphType string

	// ReliabilityModel - specifies reliability model
	ReliabilityModel string

	// Probability - specifies expression to evaluate probability of broken edge
	Probability string

	// ProtocolName - specifies protocol
	ProtocolName string

	// Experiment - specifies experiment details
	Experiment string
}

// parseArgs - parses arguments passed in command line
// TODO - at the end check with functional reqs
func (args *AppArgs) parseArgs() {
	flag.StringVar(&args.GraphFile, "graph-file", "", "read graph structure from given file")
	flag.StringVar(&args.GraphCopyFile, "graph-copy-file", "", "save copy of graph structure to file")
	flag.StringVar(&args.StatsFile, "stats-file", "", "save statistics to file")
	flag.StringVar(&args.GraphType, "graph-type", "", "provide graph-type "+
		"(path,$number_of_vertices|clique,$number_of_vertices|regular,$number_of_vertices,$degree|grid,$height,$width|hypercube,$dimension"+
		"|tree,$number_of_vertices,$degree|gridOfCliques,$height,$width,$number_of_vertices_in_clique)")
	flag.StringVar(&args.ReliabilityModel, "reliability-model", "", "specifies reliability model")
	flag.StringVar(&args.Probability, "p", "0.0", "specifies probability expression for reliability model")
	flag.StringVar(&args.ProtocolName, "protocol", "", "specifies protocol ('hll'|'minPropagation')")
	flag.StringVar(&args.Experiment, "experiment", "", "specifies experiment details "+
		"('extremaPropagation,$min,$max,$step,$repetitions'|countDistinct,$min,$max,$step,$repetitions')")
}

// InitializeAppArgs - initializes and validates arguments
func InitializeAppArgs() AppArgs {
	var args = AppArgs{}
	args.parseArgs()
	flag.Parse()

	if args.GraphFile != "" && args.GraphType != "" {
		fmt.Println("You cannot use graph file while trying to build")
	} else if args.GraphFile == "" && args.GraphType == "" && args.Experiment == "" {
		fmt.Println("You have to specify graph file or graph type")
		//TODO error
	}

	return args
}
