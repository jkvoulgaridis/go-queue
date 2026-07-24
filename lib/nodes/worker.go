package nodes

import (
	"bufio"
	"gqueue/v2/lib/network"
	"gqueue/v2/lib/messages"
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

func handleTask(task messages.TaskEnvelope, brokerAddr string) {
	log.Printf("worker task start id=%s name=%s payload=%q", task.Eid, task.Task.TaskName, task.Task.Payload)
	handler, err := network.GetTaskHandler(task.Task.TaskName)
	if err != nil {
		log.Printf("worker task handler_missing id=%s name=%s err=%v", task.Eid, task.Task.TaskName, err)
		sendNack(brokerAddr, task.Eid, err.Error())
		return
	}
	if err := handler.Run(task.Task.Payload); err != nil {
		log.Printf("worker task failed id=%s name=%s err=%v", task.Eid, task.Task.TaskName, err)
		sendNack(brokerAddr, task.Eid, err.Error())
		return
	}
	log.Printf("worker task success id=%s name=%s", task.Eid, task.Task.TaskName)
	sendAck(brokerAddr, task.Eid)
}

func workerCycle(brokerAddr string, concurrency int) error {
	conn, err := net.Dial("tcp", brokerAddr)
	if err != nil {
		return err
	}
	defer conn.Close()

	pull := messages.PullMesasges{Type: "pull", Concurrency: concurrency}
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

	env, err := network.FromCanonicalJson[messages.MultiTaskEnvelope](resp)
	if err != nil {
		log.Printf("worker pull parse_error err=%v raw=%q", err, resp)
		return err
	}

	log.Printf("worker pull received assigned=%d", len(env.Tasks))

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
	_, _ = conn.Write([]byte(network.ToCanonicalJson(messages.AckMessage{Type: "ack", Eid: taskID})))
}

func sendNack(brokerAddr, taskID, reason string) {
	conn, err := net.Dial("tcp", brokerAddr)
	if err != nil {
		log.Printf("failed to send nack for %s: %v", taskID, err)
		return
	}
	defer conn.Close()
	_, _ = conn.Write([]byte(network.ToCanonicalJson(messages.NackMessage{Type: "nack", Eid: taskID, Reason: reason})))
}
