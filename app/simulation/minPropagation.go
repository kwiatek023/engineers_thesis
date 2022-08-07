package simulation

import (
	"fmt"
	"math/rand"
)

type MinPropagation struct{}

func (MinPropagation) GetInitialData(station *Station) {
	randValue := rand.ExpFloat64()
	data := []float64{randValue}
	station.SetCurrentData(data)
}

func (MinPropagation) OnInitialize(station *Station) {
	data := station.GetCurrentData()
	station.UserDefinedVariables["min"] = data[0]
	station.Broadcast()
}

func (MinPropagation) OnDataReceive(station *Station) {
	min := station.UserDefinedVariables["min"].(float64)
	mq := station.GetMsgQueue()
	for mq.Len() > 0 {
		msg := mq.Dequeue()
		value := msg.Data[0]
		if value < min {
			min = value
		}
	}

	station.UserDefinedVariables["min"] = min
}

func (MinPropagation) OnDataPropagate(station *Station) {
	min := station.UserDefinedVariables["min"].(float64)
	currentValue := station.GetCurrentData()[0]
	if min < currentValue {
		station.SetCurrentData([]float64{min})
		// state changed, inform neighbours
		station.SynchronizedBroadcast()
	}

}

func (MinPropagation) StopCondition(station *Station) bool {
	return station.GetRoundCounter() < 6
}

func (MinPropagation) OnFinalize(station *Station) {
	fmt.Println("result of", station.GetId(),
		station.GetCurrentData()[0], "msg: ", station.GetSentMsgCounter())
}
