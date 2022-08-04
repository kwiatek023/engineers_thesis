package simulation

import (
	"app/simulationGraph"
	"math/rand"
	"sync"
)

type Station struct {
	id                   int
	communicationChannel chan float64
	currentData          float64
	nofNeighbours        int
	isActive             bool
	mutex                *sync.Mutex
}

func NewStation(id int, g *simulationGraph.GraphWrapper) *Station {
	nofNeighbours := g.GraphStructure.Degree(id)
	communicationChannel := make(chan float64, nofNeighbours)

	return &Station{id: id,
		communicationChannel: communicationChannel,
		currentData:          rand.ExpFloat64(), // probably in onInitialize
		nofNeighbours:        nofNeighbours,
		isActive:             true,
		mutex:                &sync.Mutex{}}
}
