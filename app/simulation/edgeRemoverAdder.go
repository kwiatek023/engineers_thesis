package simulation

import (
	"app/simulationGraph"
	"fmt"
	"math/rand"
)

type edgeRemoverAdder struct {
	g                       *simulationGraph.GraphWrapper
	edgeUpdateBeginChannel  chan bool
	edgeUpdateFinishChannel chan bool
}

func newEdgeRemoverAdder(g *simulationGraph.GraphWrapper, edgeUpdateBeginChannel chan bool, edgeUpdateFinishChannel chan bool) *edgeRemoverAdder {
	return &edgeRemoverAdder{g: g,
		edgeUpdateBeginChannel:  edgeUpdateBeginChannel,
		edgeUpdateFinishChannel: edgeUpdateFinishChannel}
}

func (this *edgeRemoverAdder) RunEdgeUpdating(done chan bool) {
	nofVertices := this.g.GraphStructure.Order()
	edges := this.g.GetEdges()
	relMap := this.g.GetRelMap()
	for {
		for i := 0; i < nofVertices; i++ {
			select {
			case <-this.edgeUpdateBeginChannel:
				continue
			case <-done:
				fmt.Println("updater: finish")
				return
			}
		}

		for v, e := range edges {
			for w, _ := range e {
				randVal := rand.Float64()
				p := relMap[v][w]
				q := 1 - p
				if randVal < p && this.g.GraphStructure.Edge(v, w) {
					this.g.GraphStructure.DeleteBoth(v, w)
				} else if randVal < q && !this.g.GraphStructure.Edge(v, w) {
					this.g.GraphStructure.AddBoth(v, w)
				}
			}
		}

		for i := 0; i < nofVertices; i++ {
			this.edgeUpdateFinishChannel <- true
		}
	}
}
