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

func (this *Station) RunProtocol(protocol Protocol, wg *sync.WaitGroup) {
	defer wg.Done()
	protocol.GetInitialData(this)
	// round 0
	protocol.OnInitialize(this)

	// concrete rounds
	for protocol.StopCondition(this) {
		this.manager.b.WaitAtFirstBarrier()
		close(this.communicationChannel)
		this.receiveMsgs()

		protocol.OnDataReceive(this)

		this.communicationChannel = make(chan *Pack, this.nofNeighbours)
		this.manager.b.WaitAtSecondBarrier()

		protocol.OnDataPropagate(this)
		this.roundCounter++
	}

	// sum up round
	this.manager.b.WaitAtFirstBarrier()
	protocol.OnFinalize(this)
	close(this.communicationChannel)
	this.manager.b.WaitAtSecondBarrier()
}

func (this *Station) sendMsgToStation(receiverId int) {
	packToSend := NewPack(this.currentData, this.roundCounter)
	this.manager.getStationById(receiverId).communicationChannel <- packToSend
	this.sentMsgCounter++
}

func (this *Station) receiveMsgs() {
	for msg := range this.communicationChannel {
		this.msgQueue.Enqueue(msg)
		this.receivedMsgCounter++
	}
}

func (this *Station) Broadcast() {
	this.graph.GraphStructure.Visit(this.id, func(w int, c int64) (skip bool) {
		this.sendMsgToStation(w)
		return
	})
}

func (this *Station) SynchronizedBroadcast() {
	this.graph.GraphStructure.Visit(this.id, func(w int, c int64) (skip bool) {
		this.manager.getStationById(w).mutex.Lock()
		this.sendMsgToStation(w)
		this.manager.getStationById(w).mutex.Unlock()
		return
	})
}

func (this *Station) SetCurrentData(data []float64) {
	this.currentData = data
}

func (this *Station) GetCurrentData() []float64 {
	return this.currentData
}

func (this *Station) GetMsgQueue() *MessageQueue {
	return this.msgQueue
}

func (this Station) GetId() int {
	return this.id
}

func (this Station) GetSentMsgCounter() int {
	return this.sentMsgCounter
}

func (this Station) GetRoundCounter() int {
	return this.roundCounter
}
