package simulation

import (
	"fmt"
	"math/rand"
)

type MinPropagationProtocol struct{}

func (MinPropagationProtocol) GetInitialData(station IStation) {
	randValue := rand.ExpFloat64()
	data := []float64{randValue}
	station.SetCurrentData(data)
}

func (MinPropagationProtocol) OnInitialize(station IStation) {
	data := station.GetCurrentData()
	station.SetUserDefinedVariable("min", data[0])
	station.Broadcast()
}

func (MinPropagationProtocol) OnDataReceive(station IStation) {
	min := station.GetUserDefinedVariable("min").(float64)
	mq := station.GetMsgQueue()
	for mq.Len() > 0 {
		msg := mq.Dequeue()
		value := msg.Data[0]
		if value < min {
			min = value
		}
	}

	station.SetUserDefinedVariable("min", min)
}

func (MinPropagationProtocol) OnDataPropagate(station IStation) {
	min := station.GetUserDefinedVariable("min").(float64)
	currentValue := station.GetCurrentData()[0]
	if min < currentValue {
		station.SetCurrentData([]float64{min})
		// state changed, inform neighbours
		station.SynchronizedBroadcast()
	}

}

func (MinPropagationProtocol) StopCondition(station IStation) bool {
	return station.GetRoundCounter() < station.GetGraph().GetDiameter()
}

func (MinPropagationProtocol) OnFinalize(station IStation) {
	fmt.Println("result of", station.GetId(),
		station.GetCurrentData()[0], "msg: ", station.GetSentMsgCounter())
}
