package simulation

import (
	"app/simulationGraph"
	"sync"
)

type SynchronousStation struct {
	*Station
	manager              *Manager
	communicationChannel chan *Pack
	observedValues       [][]float64
	isActive             bool
	mutex                *sync.Mutex
}

func NewSynchronousStation(manager *Manager, id int, g *simulationGraph.GraphWrapper) *SynchronousStation {
	nofNeighbours := g.GraphStructure.Degree(id)
	communicationChannel := make(chan *Pack, nofNeighbours)
	observedValues := make([][]float64, 0)

	return &SynchronousStation{NewStation(id, g),
		manager,
		communicationChannel,
		observedValues,
		true,
		&sync.Mutex{}}
}

func (this *SynchronousStation) RunProtocol(protocol Protocol, wg *sync.WaitGroup) {
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

func (this *SynchronousStation) sendMsgToStation(receiverId int) {
	packToSend := NewPack(this.currentData, this.roundCounter)
	this.manager.getStationById(receiverId).communicationChannel <- packToSend
	this.sentMsgCounter++
}

func (this *SynchronousStation) receiveMsgs() {
	for msg := range this.communicationChannel {
		this.msgQueue.Enqueue(msg)
		this.receivedMsgCounter++
	}
}

func (this *SynchronousStation) Broadcast() {
	this.graph.GraphStructure.Visit(this.id, func(w int, c int64) (skip bool) {
		this.sendMsgToStation(w)
		return
	})
}

func (this *SynchronousStation) SynchronizedBroadcast() {
	this.graph.GraphStructure.Visit(this.id, func(w int, c int64) (skip bool) {
		this.manager.getStationById(w).mutex.Lock()
		this.sendMsgToStation(w)
		this.manager.getStationById(w).mutex.Unlock()
		return
	})
}
