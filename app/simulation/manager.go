package simulation

import (
	"app/simulationGraph"
	"app/threading/barrier"
	"sync"
)

type Manager struct {
	nofStations       int
	graph             *simulationGraph.GraphWrapper
	stations          *[]*Station
	nofActiveStations int
	b                 *barrier.Barrier
}

func NewManager(nofVertices int, graph *simulationGraph.GraphWrapper) *Manager {
	stations := make([]*Station, 0)
	b := barrier.New(nofVertices)

	manager := &Manager{nofStations: nofVertices,
		stations: &stations,
		graph:    graph,
		b:        b}

	for i := 0; i < nofVertices; i++ {
		stations = append(stations, NewStation(manager, i, graph))
	}

	return manager
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

func (m Manager) getStationById(id int) *Station {
	return (*m.stations)[id]
}
