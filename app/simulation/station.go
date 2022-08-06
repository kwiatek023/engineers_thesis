package simulation

import (
	"app/protocols"
	"app/simulationGraph"
	"sync"
)

type Station struct {
	manager              Manager
	id                   int
	communicationChannel chan *Pack
	currentData          []float64
	prevData             []float64
	nofNeighbours        int
	isActive             bool
	graph                *simulationGraph.GraphWrapper
	mutex                *sync.Mutex
	sentMsgCounter       int
	receivedMsgCounter   int
	roundCounter         int
	UserDefinedVariables map[string]interface{}
}

func NewStation(manager Manager, id int, g *simulationGraph.GraphWrapper) *Station {
	nofNeighbours := g.GraphStructure.Degree(id)
	communicationChannel := make(chan *Pack, nofNeighbours)

	return &Station{manager: manager,
		id:                   id,
		communicationChannel: communicationChannel,
		currentData:          make([]float64, 0),
		nofNeighbours:        nofNeighbours,
		isActive:             true,
		graph:                g,
		mutex:                &sync.Mutex{},
		sentMsgCounter:       0,
		receivedMsgCounter:   0,
		UserDefinedVariables: make(map[string]interface{})}
}

func (s Station) RunProtocol(protocol protocols.Protocol, wg *sync.WaitGroup) {
	defer wg.Done()
	protocol.GetInitialData(s)
	// round 0
	protocol.OnInitialize(s)

	// concrete rounds
	for protocol.StopCondition(s) {
		s.manager.b.WaitAtFirstBarrier()
		close(s.communicationChannel)

		protocol.OnDataReceive(s)

		s.communicationChannel = make(chan *Pack, s.nofNeighbours)
		s.manager.b.WaitAtSecondBarrier()

		protocol.OnDataPropagate(s)
		s.roundCounter++
	}

	// sum up round
	s.manager.b.WaitAtFirstBarrier()
	protocol.OnFinalize(s)
	close(s.communicationChannel)
	s.manager.b.WaitAtSecondBarrier()
}

func (s Station) Broadcast() {
	s.graph.GraphStructure.Visit(s.id, func(w int, c int64) (skip bool) {
		//fmt.Println("sending from", vd.id, "to", w)
		packToSend := NewPack(s.currentData, s.roundCounter)
		s.manager.stations[w].communicationChannel <- packToSend
		//fmt.Println("sent from", vd.id, "to", w)
		s.sentMsgCounter++ // should I count messages here?
		return
	})
}

func (s Station) SynchronizedBroadcast() {
	s.graph.GraphStructure.Visit(s.id, func(w int, c int64) (skip bool) {
		// acquaire neighbour (if active) mutex to send value safely
		packToSend := NewPack(s.currentData, s.roundCounter)
		s.manager.stations[w].mutex.Lock()
		// if with lock maybe timeout is legal here?
		//fmt.Println("MINIMUM: sending from", vd.id, "to", w)
		s.manager.stations[w].communicationChannel <- packToSend
		s.manager.stations[w].mutex.Unlock()
		s.sentMsgCounter++ // should I count messages here?
		return
	})
}

func (s Station) SetCurrentData(data []float64) {
	s.currentData = data
}

func (s Station) GetCurrentData() []float64 {
	return s.currentData
}

func (s Station) GetCommunicationChannel() chan *Pack {
	return s.communicationChannel
}

func (s Station) GetId() int {
	return s.id
}

func (s Station) GetSentMsgCounter() int {
	return s.sentMsgCounter
}

func (s Station) GetRoundCounter() int {
	return s.roundCounter
}
