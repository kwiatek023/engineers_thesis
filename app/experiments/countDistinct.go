package experiments

import (
	"app/io"
	"app/simulation"
	"app/simulationGraph"
	"app/utils"
	"fmt"
	"strings"
)

func CountDistinctExperiment(experimentDetails string) {
	params := strings.Split(experimentDetails, ",")
	min := utils.ParseStrToPositiveInt(params[1])
	max := utils.ParseStrToPositiveInt(params[2])
	step := utils.ParseStrToPositiveInt(params[3])
	repetitions := utils.ParseStrToPositiveInt(params[4])

	for i := min; i <= max; i += step {
		g := simulationGraph.BuildGrid(i, i, "", "")
		g.SetDiameter(2 * (i - 1))

		for j := 0; j < repetitions; j++ {
			manager := simulation.NewManager("", g)
			result := manager.RunSimulation("hll")
			filepath := fmt.Sprintf("%s_%d_%d.json", "results/countDistinct/hll", i, j)
			io.SaveStatistics(filepath, result)
		}
	}
}
