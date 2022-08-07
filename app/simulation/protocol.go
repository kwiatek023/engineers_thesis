package simulation

type Protocol interface {
	GetInitialData(station *Station)
	OnInitialize(station *Station)
	OnDataReceive(station *Station)
	OnDataPropagate(station *Station)
	StopCondition(station *Station) bool
	OnFinalize(station *Station)
}
