package simulation

import (
	"app/simulationGraph"
	"math/rand"
)

type edgeRemover struct {
	g                       *simulationGraph.GraphWrapper
	edgeUpdateBeginChannel  chan bool
	edgeUpdateFinishChannel chan bool
}

func newEdgeRemover(g *simulationGraph.GraphWrapper, edgeUpdateBeginChannel chan bool, edgeUpdateFinishChannel chan bool) *edgeRemover {
	return &edgeRemover{g: g,
		edgeUpdateBeginChannel:  edgeUpdateBeginChannel,
		edgeUpdateFinishChannel: edgeUpdateFinishChannel}
}

func (this *edgeRemover) RunEdgeUpdating(done chan bool) {
	nofVertices := this.g.GraphStructure.Order()
	edges := this.g.GetEdges()
	relMap := this.g.GetRelMap()
	for {
		for i := 0; i < nofVertices; i++ {
			select {
			case <-this.edgeUpdateBeginChannel:
				continue
			case <-done:
				return
			}
		}

		for v, e := range edges {
			for w, _ := range e {
				randVal := rand.Float64()
				if randVal < relMap[v][w] && this.g.GraphStructure.Edge(v, w) {
					this.g.GraphStructure.DeleteBoth(v, w)
				}
			}
		}

		for i := 0; i < nofVertices; i++ {
			this.edgeUpdateFinishChannel <- true
		}
	}
}
