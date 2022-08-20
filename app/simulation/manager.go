package simulation

import (
	"app/config"
	"app/io"
	"app/simulationGraph"
	"app/threading/barrier"
	"sync"
)

type Manager struct {
	nofStations       int
	graph             *simulationGraph.GraphWrapper
	stations          *[]*SynchronousStation
	nofActiveStations int
	b                 *barrier.Barrier
}

func NewManager(args config.AppArgs) *Manager {
	graph := prepareSimulationGraph(args)

	nofVertices := graph.GraphStructure.Order()
	stations := make([]*SynchronousStation, 0)
	b := barrier.New(nofVertices)

	manager := &Manager{nofStations: nofVertices,
		stations: &stations,
		graph:    graph,
		b:        b}

	for i := 0; i < nofVertices; i++ {
		stations = append(stations, NewSynchronousStation(manager, i, graph))
	}

	return manager
}

func prepareSimulationGraph(args config.AppArgs) *simulationGraph.GraphWrapper {
	var g *simulationGraph.GraphWrapper
	if args.GraphFile != "" {
		conf := io.ReadGraphFromFile(args.GraphFile)
		g = simulationGraph.BuildGraphFromConfig(conf)
	} else {
		g = simulationGraph.BuildGraphFromType(args)
	}

	return g
}

func (m Manager) RunSimulation() {
	p := MinPropagationProtocol{}
	var wg sync.WaitGroup
	wg.Add(len(*m.stations))

	for _, s := range *m.stations {
		go s.RunProtocol(p, &wg)
	}

	wg.Wait()
	m.b.Close()
}

func (m Manager) getStationById(id int) *SynchronousStation {
	return (*m.stations)[id]
}
