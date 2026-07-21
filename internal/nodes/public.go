package nodes

import (
	"fmt"
	"gqueue/v2/internal/execution"
	"gqueue/v2/internal/network"
	"log"
	"math"
	"net"
	"os"
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

	msg := network.TaskSubmission{
		Type: "submit",
		TaskToExec: network.TaskExecution{
			TaskName: taskName,
			Payload:  payload,
		},
	}

	if _, err := conn.Write([]byte(network.ToCanonicalJson(msg))); err != nil {
		return err
	}

	return nil
}

func Run(
	mode string,
	brokerAddr string,
	submitCount int,
	problematicPercent float64,
) {

	if mode == "submit" {
		addr := brokerAddr
		if addr == "" {
			addr = os.Getenv("BROKER_ADDR")
		}
		if addr == "" {
			addr = "localhost:18018"
		}

		if submitCount < 0 {
			log.Fatal("count must be >= 0")
		}
		if problematicPercent < 0 || problematicPercent > 100 {
			log.Fatal("problematic-percent must be between 0 and 100")
		}

		problematicCount := int(math.Round(float64(submitCount) * (problematicPercent / 100.0)))
		if problematicCount > submitCount {
			problematicCount = submitCount
		}
		normalCount := submitCount - problematicCount

		for i := 0; i < normalCount; i++ {
			payload := fmt.Sprintf("q2-finance-report-%d", i+1)
			if err := SubmitTask(addr, "submit_report", payload); err != nil {
				log.Fatal(err)
			}
		}
		for i := 0; i < problematicCount; i++ {
			payload := fmt.Sprintf("payload-%d", i+1)
			if err := SubmitTask(addr, "problematic", payload); err != nil {
				log.Fatal(err)
			}
		}

		log.Printf("submitted tasks to broker=%s total=%d normal=%d problematic=%d problematic_percent=%.2f", addr, submitCount, normalCount, problematicCount, problematicPercent)
		return
	}

	nodeType := os.Getenv("NODE_TYPE")
	if mode != "" {
		nodeType = mode
	}

	switch nodeType {
	case "master", "broker":
		MasterSpinup()
	case "worker":
		WorkerSpinUp()
	default:
		log.Println("NODE_TYPE not set; running broker by default")
		MasterSpinup()
	}
}
