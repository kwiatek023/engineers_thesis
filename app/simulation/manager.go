package simulation

import (
	"app/simulationGraph"
	"app/threading/barrier"
)

type Manager struct {
	nofStations       int
	graph             *simulationGraph.GraphWrapper
	stations          []*Station
	nofActiveStations int
	nofRounds         int
	b                 *barrier.Barrier
}

func NewManager(nofVertices int, nofRounds int, graph *simulationGraph.GraphWrapper) *Manager {
	var stations []*Station

	for i := 0; i < nofVertices; i++ {
		stations = append(stations, NewStation(i, graph))
	}

	b := barrier.New(nofVertices)

	return &Manager{nofStations: nofVertices,
		stations:  stations,
		graph:     graph,
		b:         b,
		nofRounds: nofRounds}
}

func (m Manager) RunSimulation() {

}
