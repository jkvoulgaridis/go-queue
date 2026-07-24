package network

import (
	"bufio"
	"gqueue/v2/lib/execution"
	"gqueue/v2/lib/messages"
	"gqueue/v2/lib/singletons"
	"io"
	"log"
	"net"
	"strings"
)

func RegisterTask(name string, handler execution.Task) error {
	return singletons.TaskRegistry.Register(name, handler)
}

func GetTaskHandler(name string) (execution.Task, error) {
	return singletons.TaskRegistry.GetHandler(name)
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
		log.Printf("broker recv raw=%q", msg)

		msgStruct, err := FromCanonicalJson[messages.Message](msg)
		if err != nil {
			log.Printf("Error parsing message envelope: %v", err)
			continue
		}

		switch msgStruct.Type {
		case "submit":
			HandleSubmitTask(msg)
			if _, err := writer.WriteString(ToCanonicalJson(messages.Message{Type: "ok"})); err != nil {
				log.Printf("Error writing submit response: %v", err)
				return
			}
			if err := writer.Flush(); err != nil {
				log.Printf("Error flushing submit response: %v", err)
				return
			}
		case "pull":
			assignPayload := HandlePullTask(msg)
			if assignPayload == "" {
				log.Printf("broker pull empty_assign_payload; returning empty task list")
				assignPayload = ToCanonicalJson(messages.MultiTaskEnvelope{Type: "assign", Tasks: []messages.TaskEnvelope{}})
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
			ackMsg, err := FromCanonicalJson[messages.AckMessage](msg)
			if err != nil {
				log.Printf("Invalid ack message: %v", err)
				continue
			}
			HandleAckMessage(ackMsg)
		case "nack":
			nackMsg, err := FromCanonicalJson[messages.NackMessage](msg)
			if err != nil {
				log.Printf("Invalid nack message: %v", err)
				continue
			}
			HandleNackkMessage(nackMsg)
			log.Printf("Got NACK for task id: %s reason: %s", nackMsg.Eid, nackMsg.Reason)
		default:
			log.Printf("Unknown message type: %s", msgStruct.Type)
		}
	}
}

func GetMetrics() string {
	taskQueue := singletons.TaskQueue
	return ToCanonicalJson(taskQueue.GetMetrics())
}
