package examples

import (
	"gqueue/v2/lib"
)

func SetupTasks() {
	err := lib.RegisterTasks([]lib.TaskDef{
		{Name: "submit_report", Handler: SubmitReportHandler{}},
		{Name: "problematic", Handler: ProblematicHandler{}},
	})
	if err != nil {
		return
	}
}
