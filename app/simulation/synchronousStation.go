package simulation

import (
	"app/simulationGraph"
	"github.com/DmitriyVTitov/size"
	"go/types"
	"sync"
)

// SynchronousStation - station used for implementing synchronous protocols
type SynchronousStation struct {
	*Station
	manager              *Manager
	communicationChannel chan *Pack
	mutex                *sync.Mutex
	maxQueueSize         int
}

func NewSynchronousStation(manager *Manager, id int, g *simulationGraph.GraphWrapper) *SynchronousStation {
	nofNeighbours := g.GraphStructure.Degree(id)
	communicationChannel := make(chan *Pack, nofNeighbours)

	return &SynchronousStation{NewStation(id, g),
		manager,
		communicationChannel,
		&sync.Mutex{},
		0}
}

// RunProtocol - runs protocol
func (this *SynchronousStation) RunProtocol(protocol Protocol, wg *sync.WaitGroup, reliabilityModel string, updateBegin chan bool,
	updateFinish chan bool) {
	defer wg.Done()
	protocol.GetInitialData(this)
	// round 0
	protocol.OnInitialize(this)

	// concrete rounds
	for protocol.StopCondition(this) {
		if reliabilityModel != "" {
			this.waitForUpdate(updateBegin, updateFinish)
		}

		this.manager.b.WaitAtFirstBarrier()
		close(this.communicationChannel)
		this.receiveMsgs()
		this.updateMaxQueueSizeIfNecessary()
		protocol.OnDataReceive(this)

		this.communicationChannel = make(chan *Pack, this.nofNeighbours)
		this.manager.b.WaitAtSecondBarrier()

		protocol.OnDataPropagate(this)
		this.RoundCounter++
	}

	// sum up round
	this.manager.b.WaitAtFirstBarrier()
	protocol.OnFinalize(this)
	close(this.communicationChannel)
	this.ExactResult = protocol.CalculateStationExactResult(this)
	this.MemoryCounter += size.Of(this.userDefinedVariables) / size.Of(types.Float64)
	this.MemoryCounter += this.maxQueueSize
	this.MemoryCounter += len(this.currentData)
	if len(this.observedValues) > 0 {
		this.MemoryCounter += len(this.observedValues) * len(this.observedValues[0])
	}
	this.manager.b.WaitAtSecondBarrier()
}

func (this *SynchronousStation) waitForUpdate(updateBegin chan bool, updateFinish chan bool) {
	updateBegin <- true
	<-updateFinish
}

func (this *SynchronousStation) updateMaxQueueSizeIfNecessary() {
	if this.msgQueue.Len() > this.maxQueueSize {
		this.maxQueueSize = this.msgQueue.Len()
	}
}

func (this *SynchronousStation) sendMsgToStation(receiverId int) {
	packToSend := NewPack(this.currentData, this.RoundCounter)
	s := this.manager.getStationById(receiverId).(*SynchronousStation)
	s.communicationChannel <- packToSend
	this.SentMsgCounter += len(this.currentData)
}

func (this *SynchronousStation) receiveMsgs() {
	for msg := range this.communicationChannel {
		this.historicalDataForStats = append(this.historicalDataForStats, msg.Data)
		this.msgQueue.Enqueue(msg)
		this.ReceivedMsgCounter += len(msg.Data)
	}
}

// Broadcast - function used for broadcasting information to neighbours
func (this *SynchronousStation) Broadcast() {
	this.graph.GraphStructure.Visit(this.id, func(w int, c int64) (skip bool) {
		this.sendMsgToStation(w)
		return
	})
}

// SynchronizedBroadcast - thread-safe function used for broadcasting information to neighbours
func (this *SynchronousStation) SynchronizedBroadcast() {
	this.graph.GraphStructure.Visit(this.id, func(w int, c int64) (skip bool) {
		s := this.manager.getStationById(w).(*SynchronousStation)
		s.mutex.Lock()
		this.sendMsgToStation(w)
		s.mutex.Unlock()
		return
	})
}

func (this *SynchronousStation) GetStation() Station {
	return *this.Station
}
