package protocols

import "app/simulation"

type Protocol interface {
	GetInitialData(station simulation.Station)
	OnInitialize(station simulation.Station)
	OnDataReceive(station simulation.Station)
	OnDataPropagate(station simulation.Station)
	StopCondition(station simulation.Station) bool
	OnFinalize(station simulation.Station)
}
