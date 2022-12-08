package simulation

// Protocol - an interface used for creating custom protocols
type Protocol interface {
	// GetInitialData - phase in which each station generates initial data
	GetInitialData(station IStation)
	// OnInitialize - phase in which round 0 is performed (setup for next protocol rounds)
	OnInitialize(station IStation)
	// OnDataReceive - in this phase stations read messages (if messages are present)
	OnDataReceive(station IStation)
	// OnDataPropagate - in this phase stations propagate some information
	OnDataPropagate(station IStation)
	// StopCondition - stop condition defines when protocol is finished
	StopCondition(station IStation) bool
	// OnFinalize - phase in which wrap-up round is performed
	OnFinalize(station IStation)
	// CalculateStationExactResult - function used to calculate exact station result
	CalculateStationExactResult(station IStation) float64
	// CalculateGlobalExactResult - function used to calculate global result using all stations
	CalculateGlobalExactResult(stations *[]IStation) float64
}
