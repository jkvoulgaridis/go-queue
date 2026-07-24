package lib

import (
	"gqueue/v2/lib/execution"
	"gqueue/v2/lib/network"
	"gqueue/v2/lib/nodes"
	"net"
	"gqueue/v2/lib/messages"
)

type TaskDef struct {
	Name    string
	Handler execution.Task
}

func RegisterTasks(tasks []TaskDef) error {
	for _, task := range tasks {
		if err := network.RegisterTask(task.Name, task.Handler); err != nil {
			return err
		}
	}
	return nil
}

func SubmitTask(brokerAddr string, taskName string, payload string) error {
	conn, err := net.Dial("tcp", brokerAddr)
	if err != nil {
		return err
	}
	defer conn.Close()

	msg := messages.TaskSubmission{
		Type: "submit",
		TaskToExec: messages.TaskExecution{
			TaskName: taskName,
			Payload:  payload,
		},
	}

	if _, err := conn.Write([]byte(network.ToCanonicalJson(msg))); err != nil {
		return err
	}

	return nil
}

func RunBroker() {
	nodes.MasterSpinup()
}

func RunWorker() {
	nodes.WorkerSpinUp()
}
