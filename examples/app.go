package examples

import (
	"gqueue/v2/internal/nodes"
)

func App(mode string, brokerAddr string, submitCount int, problematicPercent float64) {
	err := nodes.RegisterTasks([]nodes.TaskDef{
		{Name: "submit_report", Handler: SubmitReportHandler{}},
		{Name: "problematic", Handler: ProblematicHandler{}},
	})
	if err != nil {
		return
	}
	nodes.Run(mode, brokerAddr, submitCount, problematicPercent)
}
