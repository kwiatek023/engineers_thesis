package protocols

import (
	"app/simulation"
	"fmt"
	"math/rand"
)

type MinPropagation struct{}

func (MinPropagation) GetInitialData(station simulation.Station) {
	randValue := rand.ExpFloat64()
	data := []float64{randValue}
	station.SetCurrentData(data)
}

func (MinPropagation) OnInitialize(station simulation.Station) {
	data := station.GetCurrentData()
	station.UserDefinedVariables["min"] = data[0]
	station.Broadcast()
}

func (MinPropagation) OnDataReceive(station simulation.Station) {
	min := station.UserDefinedVariables["min"].(float64)
	for msg := range station.GetCommunicationChannel() {
		value := msg.Data[0]
		if value < min {
			station.UserDefinedVariables["min"] = msg.Data[0]
		}
	}
}

func (MinPropagation) OnDataPropagate(station simulation.Station) {
	min := station.UserDefinedVariables["min"].(float64)
	currentValue := station.GetCurrentData()[0]
	if min < currentValue {
		station.SetCurrentData([]float64{currentValue})
		// state changed, inform neighbours
		station.SynchronizedBroadcast()
	}

}

func (MinPropagation) StopCondition(station simulation.Station) bool {
	return station.GetRoundCounter() < 6
}

func (MinPropagation) OnFinalize(station simulation.Station) {
	fmt.Println("result of", station.GetId(),
		station.GetCurrentData()[0], "msg: ", station.GetSentMsgCounter())
}
