package simulation

import (
	"app/simulationGraph"
	"app/threading/barrier"
	"github.com/montanaflynn/stats"
	"sync"
)

type Manager struct {
	nofStations       int
	graph             *simulationGraph.GraphWrapper
	stations          *[]IStation
	nofActiveStations int
	reliabilityModel  string
	b                 *barrier.Barrier
}

func NewManager(reliabilityModel string, graph *simulationGraph.GraphWrapper) *Manager {
	nofVertices := graph.GraphStructure.Order()
	stations := make([]IStation, 0)
	b := barrier.New(nofVertices)

	manager := &Manager{nofStations: nofVertices,
		stations:         &stations,
		graph:            graph,
		reliabilityModel: reliabilityModel,
		b:                b}

	for i := 0; i < nofVertices; i++ {
		stations = append(stations, NewSynchronousStation(manager, i, graph))
	}

	return manager
}

func (m Manager) RunSimulation(protocolName string) JsonStatsStructure {
	var wg sync.WaitGroup
	p := m.mapNameToProtocol(protocolName)
	wg.Add(len(*m.stations))
	updateBeginChannel := make(chan bool, 1)
	updateFinishChannel := make(chan bool, 1)
	done := make(chan bool)
	updater := m.getReliabilityModel(m.reliabilityModel, updateBeginChannel, updateFinishChannel)

	if m.reliabilityModel != "" {
		go updater.RunEdgeUpdating(done)
	}

	for _, s := range *m.stations {
		go s.RunProtocol(p, &wg, m.reliabilityModel, updateBeginChannel, updateFinishChannel)
	}

	wg.Wait()
	if m.reliabilityModel != "" {
		done <- true
	}
	m.b.Close()
	close(updateBeginChannel)
	close(updateFinishChannel)
	return m.makeStatsSummary(p)
}

func (m Manager) getStationById(id int) IStation {
	return (*m.stations)[id]
}

func (m Manager) getReliabilityModel(reliabilityModel string, updateBeginChannel, updateFinishChannel chan bool) IEdgeUpdater {
	if reliabilityModel == "edge-remover" {
		return newEdgeRemover(m.graph, updateBeginChannel, updateFinishChannel)
	} else if reliabilityModel == "edge-remover-adder" {
		return newEdgeRemoverAdder(m.graph, updateBeginChannel, updateFinishChannel)
	}
	return nil
}

func (m Manager) makeStatsSummary(p Protocol) JsonStatsStructure {
	exactResult := p.CalculateGlobalExactResult(m.stations)
	stations := make([]Station, 0)
	msgsSentStats := make([]float64, 0)
	msgsReceivedStats := make([]float64, 0)
	roundsStats := make([]float64, 0)
	memoryStats := make([]float64, 0)

	for _, station := range *m.stations {
		msgsSentStats = append(msgsSentStats, float64(station.GetSentMsgCounter()))
		msgsReceivedStats = append(msgsReceivedStats, float64(station.GetReceivedMsgCounter()))
		roundsStats = append(roundsStats, float64(station.GetRoundCounter()))
		memoryStats = append(memoryStats, float64(station.GetMemoryCounter()))
		stations = append(stations, station.GetStation())
	}

	nofRounds, _ := stats.Max(roundsStats)
	maxReceivedMsgs, _ := stats.Max(msgsReceivedStats)
	minReceivedMsgs, _ := stats.Min(msgsReceivedStats)
	allReceivedMsgs, _ := stats.Sum(msgsReceivedStats)
	avgReceivedMsgs, _ := stats.Mean(msgsReceivedStats)
	stddevReceivedMsgs, _ := stats.StandardDeviation(msgsReceivedStats)

	maxSentMsgs, _ := stats.Max(msgsSentStats)
	minSentMsgs, _ := stats.Min(msgsSentStats)
	allSentMsgs, _ := stats.Sum(msgsSentStats)
	avgSentMsgs, _ := stats.Mean(msgsSentStats)
	stddevSentMsgs, _ := stats.StandardDeviation(msgsSentStats)

	allMemory, _ := stats.Sum(memoryStats)
	maxMemory, _ := stats.Max(memoryStats)
	minMemory, _ := stats.Min(memoryStats)
	avgMemory, _ := stats.Mean(memoryStats)
	stddevMemory, _ := stats.StandardDeviation(memoryStats)

	statistics := JsonStatsStructure{
		Size:               m.graph.GraphStructure.Order(),
		Result:             exactResult,
		NofRounds:          int(nofRounds),
		MaxReceivedMsgs:    int(maxReceivedMsgs),
		MinReceivedMsgs:    int(minReceivedMsgs),
		AllReceivedMsgs:    int(allReceivedMsgs),
		AvgReceivedMsgs:    avgReceivedMsgs,
		StddevReceivedMsgs: stddevReceivedMsgs,
		MaxSentMsgs:        int(maxSentMsgs),
		MinSentMsgs:        int(minSentMsgs),
		AllSentMsgs:        int(allSentMsgs),
		AvgSentMsgs:        avgSentMsgs,
		StddevSentMsgs:     stddevSentMsgs,
		AllMemory:          int(allMemory),
		MaxMemory:          int(maxMemory),
		MinMemory:          int(minMemory),
		AvgMemory:          avgMemory,
		StddevMemory:       stddevMemory,
		Stations:           stations,
	}

	return statistics
}
