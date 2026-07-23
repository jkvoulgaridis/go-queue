package nodes

import (
	"bufio"
	"gqueue/v2/lib/network"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"time"
)

func WorkerSpinUp() {
	brokerAddr := os.Getenv("BROKER_ADDR")
	if brokerAddr == "" {
		brokerAddr = "master:18018"
	}

	concurrency := 1
	if c := os.Getenv("WORKER_CONCURRENCY"); c != "" {
		if parsed, err := strconv.Atoi(c); err == nil && parsed > 0 {
			concurrency = parsed
		}
	}

	pollInterval := 2 * time.Second
	if ms := os.Getenv("WORKER_POLL_MS"); ms != "" {
		if parsed, err := strconv.Atoi(ms); err == nil && parsed > 0 {
			pollInterval = time.Duration(parsed) * time.Millisecond
		}
	}

	log.Printf("worker started broker=%s concurrency=%d", brokerAddr, concurrency)

	for {
		if err := workerCycle(brokerAddr, concurrency); err != nil {
			log.Printf("worker cycle error: %v", err)
		}
		time.Sleep(pollInterval)
	}
}

func handleTask(task network.TaskEnvelope, brokerAddr string) {
	log.Println("Started running task")
	handler, err := network.GetTaskHandler(task.Task.TaskName)
	if err != nil {
		sendNack(brokerAddr, task.Eid, err.Error())
		return
	}
	if err := handler.Run(task.Task.Payload); err != nil {
		sendNack(brokerAddr, task.Eid, err.Error())
		return
	}
	sendAck(brokerAddr, task.Eid)
}

func workerCycle(brokerAddr string, concurrency int) error {
	conn, err := net.Dial("tcp", brokerAddr)
	if err != nil {
		return err
	}
	defer conn.Close()

	pull := network.PullMesasges{Type: "pull", Concurrency: concurrency}
	if _, err := conn.Write([]byte(network.ToCanonicalJson(pull))); err != nil {
		return err
	}

	reader := bufio.NewReader(conn)
	resp, err := reader.ReadString('\n')
	if err != nil {
		if err == io.EOF {
			return nil
		}
		return err
	}

	env, err := network.FromCanonicalJson[network.MultiTaskEnvelope](resp)
	if err != nil {
		return err
	}

	for _, task := range env.Tasks {
		go handleTask(task, brokerAddr)
	}

	return nil
}

func sendAck(brokerAddr, taskID string) {
	conn, err := net.Dial("tcp", brokerAddr)
	if err != nil {
		log.Printf("failed to send ack for %s: %v", taskID, err)
		return
	}
	defer conn.Close()
	_, _ = conn.Write([]byte(network.ToCanonicalJson(network.AckMessage{Type: "ack", Eid: taskID})))
}

func sendNack(brokerAddr, taskID, reason string) {
	conn, err := net.Dial("tcp", brokerAddr)
	if err != nil {
		log.Printf("failed to send nack for %s: %v", taskID, err)
		return
	}
	defer conn.Close()
	_, _ = conn.Write([]byte(network.ToCanonicalJson(network.NackMessage{Type: "nack", Eid: taskID, Reason: reason})))
}
