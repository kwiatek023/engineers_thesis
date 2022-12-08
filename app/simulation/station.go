package simulation

import (
	"app/simulationGraph"
	"sync"
)

// IStation - interface used for station module
type IStation interface {
	// RunProtocol - runs protocol
	RunProtocol(protocol Protocol, wg *sync.WaitGroup, reliabilityModel string, updateBegin chan bool,
		updateFinish chan bool)
	// GetId - returns station id
	GetId() int
	// Broadcast - sends msg to station neighbours
	Broadcast()
	// SynchronizedBroadcast - sends msg to station neighbours (threadsafe)
	SynchronizedBroadcast()
	// SetCurrentData - sets current vector data in station
	SetCurrentData(data []float64)
	// GetCurrentData - returns current vector data in station
	GetCurrentData() []float64
	// GetMsgQueue - returns station's message queue
	GetMsgQueue() *MessageQueue
	// GetSentMsgCounter - returns sent message counter
	GetSentMsgCounter() int
	// GetReceivedMsgCounter - returns received message counter
	GetReceivedMsgCounter() int
	// GetMemoryCounter - returns memory counter
	GetMemoryCounter() int
	// GetRoundCounter - returns round counter
	GetRoundCounter() int
	// SetUserDefinedVariable - sets user defined variable (key - name of variable, value - value of variable)
	SetUserDefinedVariable(key string, value interface{})
	// GetUserDefinedVariable - gets user defined variable by name
	GetUserDefinedVariable(key string) interface{}
	// GetGraph - returns graph topology
	GetGraph() *simulationGraph.GraphWrapper
	// GetHistoricalDataForStats - returns historical data for statistics
	GetHistoricalDataForStats() [][]float64
	// SetResult - sets station result
	SetResult(result float64)
	// ObserveValue - defines behaviour when station observes a value
	ObserveValue(value []float64)
	// GetObservedValues - returns observed values
	GetObservedValues() [][]float64
	// GetStation - returns station object
	GetStation() Station
}

type Station struct {
	IStation               `json:",omitempty"`
	id                     int `json:"id"`
	nofNeighbours          int
	msgQueue               *MessageQueue
	currentData            []float64
	historicalDataForStats [][]float64
	observedValues         [][]float64
	graph                  *simulationGraph.GraphWrapper
	SentMsgCounter         int `json:"sent_msgs"`
	ReceivedMsgCounter     int `json:"received_msgs"`
	RoundCounter           int `json:"nof_rounds"`
	userDefinedVariables   map[string]interface{}
	Result                 float64 `json:"result"`
	ExactResult            float64 `json:"exact_result"`
	MemoryCounter          int     `json:"memory"`
}

func NewStation(id int, graph *simulationGraph.GraphWrapper) *Station {
	nofNeighbours := graph.GraphStructure.Degree(id)
	return &Station{id: id,
		nofNeighbours:          nofNeighbours,
		msgQueue:               NewMessageQueue(),
		currentData:            make([]float64, 0),
		historicalDataForStats: make([][]float64, 0),
		observedValues:         make([][]float64, 0),
		graph:                  graph,
		SentMsgCounter:         0,
		ReceivedMsgCounter:     0,
		RoundCounter:           0,
		userDefinedVariables:   make(map[string]interface{}),
		MemoryCounter:          0}
}

func (this *Station) SetCurrentData(data []float64) {
	this.historicalDataForStats = append(this.historicalDataForStats, data)
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
	return this.SentMsgCounter
}

func (this *Station) GetReceivedMsgCounter() int {
	return this.ReceivedMsgCounter
}

func (this *Station) GetMemoryCounter() int {
	return this.MemoryCounter
}

func (this *Station) GetRoundCounter() int {
	return this.RoundCounter
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

func (this *Station) GetHistoricalDataForStats() [][]float64 {
	return this.historicalDataForStats
}

func (this *Station) SetResult(result float64) {
	this.Result = result
}

func (this *Station) ObserveValue(value []float64) {
	this.observedValues = append(this.observedValues, value)
}

func (this *Station) GetObservedValues() [][]float64 {
	return this.observedValues
}
