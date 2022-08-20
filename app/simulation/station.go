package simulation

import (
	"app/simulationGraph"
)

type IStation interface {
	GetId() int
	Broadcast()
	SynchronizedBroadcast()
	SetCurrentData(data []float64)
	GetCurrentData() []float64
	GetMsgQueue() *MessageQueue
	GetSentMsgCounter() int
	GetRoundCounter() int
	SetUserDefinedVariable(key string, value interface{})
	GetUserDefinedVariable(key string) interface{}
	GetGraph() *simulationGraph.GraphWrapper
}

type Station struct {
	IStation
	id                   int
	nofNeighbours        int
	msgQueue             *MessageQueue
	currentData          []float64
	graph                *simulationGraph.GraphWrapper
	sentMsgCounter       int
	receivedMsgCounter   int
	roundCounter         int
	userDefinedVariables map[string]interface{}
}

func NewStation(id int, graph *simulationGraph.GraphWrapper) *Station {
	nofNeighbours := graph.GraphStructure.Degree(id)
	return &Station{id: id,
		nofNeighbours:        nofNeighbours,
		msgQueue:             NewMessageQueue(),
		currentData:          make([]float64, 0),
		graph:                graph,
		sentMsgCounter:       0,
		receivedMsgCounter:   0,
		roundCounter:         0,
		userDefinedVariables: make(map[string]interface{})}
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

func (this *Station) GetId() int {
	return this.id
}

func (this *Station) GetSentMsgCounter() int {
	return this.sentMsgCounter
}

func (this *Station) GetRoundCounter() int {
	return this.roundCounter
}

func (this *Station) SetUserDefinedVariable(key string, value interface{}) {
	this.userDefinedVariables[key] = value
}

func (this *Station) GetUserDefinedVariable(key string) interface{} {
	return this.userDefinedVariables[key]
}

func (this *Station) GetGraph() *simulationGraph.GraphWrapper {
	return this.graph
}
