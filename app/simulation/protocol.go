package simulation

type Protocol interface {
	GetInitialData(station IStation)
	OnInitialize(station IStation)
	OnDataReceive(station IStation)
	OnDataPropagate(station IStation)
	StopCondition(station IStation) bool
	OnFinalize(station IStation)
	CalculateStationExactResult(station IStation) float64
	CalculateGlobalExactResult(stations *[]IStation) float64
}
