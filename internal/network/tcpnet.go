package network

import (
	"bufio"
	"gqueue/v2/internal/execution"
	"io"
	"log"
	"net"
	"strings"
)

var (
	taskRegistry = execution.NewRegistry()
	taskQueue    = execution.NewTaskQueue()
)


func RegisterTask(name string, handler execution.Task) error {
	return taskRegistry.Register(name, handler)
}

func GetTaskHandler(name string) (execution.Task, error) {
	return taskRegistry.GetHandler(name)
}

func HandleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				log.Printf("Error reading connection: %v", err)
			}
			return
		}

		msg = strings.TrimSpace(msg)
		if msg == "" {
			continue
		}

		msgStruct, err := FromCanonicalJson[Message](msg)
		if err != nil {
			log.Printf("Error parsing message envelope: %v", err)
			continue
		}

		switch msgStruct.Type {
		case "submit":
			HandleSubmitTask(msg, taskQueue, taskRegistry)
			if _, err := writer.WriteString(ToCanonicalJson(Message{Type: "ok"})); err != nil {
				log.Printf("Error writing submit response: %v", err)
				return
			}
			if err := writer.Flush(); err != nil {
				log.Printf("Error flushing submit response: %v", err)
				return
			}
		case "pull":
			assignPayload := HandlePullTask(msg, taskQueue)
			if assignPayload == "" {
				assignPayload = ToCanonicalJson(MultiTaskEnvelope{Type: "assign", Tasks: []TaskEnvelope{}})
			}
			if _, err := writer.WriteString(assignPayload); err != nil {
				log.Printf("Error writing pull response: %v", err)
				return
			}
			if err := writer.Flush(); err != nil {
				log.Printf("Error flushing pull response: %v", err)
				return
			}
		case "assign":
			log.Printf("broker received unexpected assign message")
		case "ack":
			ackMsg, err := FromCanonicalJson[AckMessage](msg)
			if err != nil {
				log.Printf("Invalid ack message: %v", err)
				continue
			}
			log.Printf("Got ACK for task id: %s", ackMsg.Eid)
		case "nack":
			nackMsg, err := FromCanonicalJson[NackMessage](msg)
			if err != nil {
				log.Printf("Invalid nack message: %v", err)
				continue
			}
			log.Printf("Got NACK for task id: %s reason: %s", nackMsg.Eid, nackMsg.Reason)
		default:
			log.Printf("Unknown message type: %s", msgStruct.Type)
		}
	}
}
