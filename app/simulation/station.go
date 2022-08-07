package simulation

import (
	"app/simulationGraph"
	"sync"
)

type Station struct {
	manager              *Manager
	id                   int
	communicationChannel chan *Pack
	msgQueue             *MessageQueue
	observedValues       [][]float64
	currentData          []float64
	nofNeighbours        int
	isActive             bool
	graph                *simulationGraph.GraphWrapper
	mutex                *sync.Mutex
	sentMsgCounter       int
	receivedMsgCounter   int
	roundCounter         int
	UserDefinedVariables map[string]interface{}
}

func NewStation(manager *Manager, id int, g *simulationGraph.GraphWrapper) *Station {
	nofNeighbours := g.GraphStructure.Degree(id)
	communicationChannel := make(chan *Pack, nofNeighbours)

	return &Station{manager: manager,
		id:                   id,
		communicationChannel: communicationChannel,
		msgQueue:             NewMessageQueue(),
		currentData:          make([]float64, 0),
		nofNeighbours:        nofNeighbours,
		isActive:             true,
		graph:                g,
		mutex:                &sync.Mutex{},
		sentMsgCounter:       0,
		receivedMsgCounter:   0,
		UserDefinedVariables: make(map[string]interface{})}
}

func (s *Station) RunProtocol(protocol Protocol, wg *sync.WaitGroup) {
	defer wg.Done()
	protocol.GetInitialData(s)
	if s.id == 0 {
		s.currentData[0] = 0
	}
	// round 0
	protocol.OnInitialize(s)

	// concrete rounds
	for protocol.StopCondition(s) {
		s.manager.b.WaitAtFirstBarrier()
		close(s.communicationChannel)
		s.receiveMsgs()

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

func (s *Station) sendMsgToStation(receiverId int) {
	packToSend := NewPack(s.currentData, s.roundCounter)
	s.manager.getStationById(receiverId).communicationChannel <- packToSend
	s.sentMsgCounter++
}

func (s *Station) receiveMsgs() {
	for msg := range s.communicationChannel {
		s.msgQueue.Enqueue(msg)
		s.receivedMsgCounter++
	}
}

func (s *Station) Broadcast() {
	s.graph.GraphStructure.Visit(s.id, func(w int, c int64) (skip bool) {
		s.sendMsgToStation(w)
		return
	})
}

func (s *Station) SynchronizedBroadcast() {
	s.graph.GraphStructure.Visit(s.id, func(w int, c int64) (skip bool) {
		s.manager.getStationById(w).mutex.Lock()
		s.sendMsgToStation(w)
		s.manager.getStationById(w).mutex.Unlock()
		return
	})
}

func (s *Station) SetCurrentData(data []float64) {
	s.currentData = data
}

func (s *Station) GetCurrentData() []float64 {
	return s.currentData
}

func (s *Station) GetMsgQueue() *MessageQueue {
	return s.msgQueue
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
