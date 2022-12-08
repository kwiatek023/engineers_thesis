package simulation

// Pack - structure for holding message data
type Pack struct {
	Data        []float64
	RoundNumber int
}

func NewPack(data []float64, roundNumber int) *Pack {
	return &Pack{Data: data, RoundNumber: roundNumber}
}
