package examples

import (
	"fmt"
	"gqueue/v2/internal/execution"
	"gqueue/v2/internal/network"
	"log"
)

func ExampleFlow() error {
	registry := execution.NewRegistry()
	queue := execution.NewTaskQueue()

	if err := registry.Register("submit_report", SubmitReportHandler{}); err != nil {
		return err
	}
	if err := registry.Register("problematic", ProblematicHandler{}); err != nil {
		return err
	}

	submit1 := network.ToCanonicalJson(network.TaskSubmission{
		Type: "submit",
		TaskToExec: network.TaskExecution{
			TaskName: "submit_report",
			Payload:  "q2-finance-report",
		},
	})
	network.HandleSubmitTask(submit1, queue, registry)

	submit2 := network.ToCanonicalJson(network.TaskSubmission{
		Type: "submit",
		TaskToExec: network.TaskExecution{
			TaskName: "problematic",
			Payload:  "payload",
		},
	})
	network.HandleSubmitTask(submit2, queue, registry)

	pull := network.ToCanonicalJson(network.PullMesasges{Type: "pull", Concurrency: 2})
	assignPayload := network.HandlePullTask(pull, queue)
	if assignPayload == "" {
		return fmt.Errorf("no assign payload generated")
	}

	log.Printf("assign payload: %s", assignPayload)
	network.HandleAssignTask(assignPayload, registry)

	return nil
}
