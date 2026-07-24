package network

import (
	"bufio"
	"gqueue/v2/lib/execution"
	"gqueue/v2/lib/singletons"
	"io"
	"log"
	"net"
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
		rsp := HanldeIncomingMessage(msg)
		if rsp == "" {
			continue
		}
		if _, err := writer.WriteString(rsp); err != nil {
			log.Printf("Error writing submit response: %v", err)
			return
		}
		if err := writer.Flush(); err != nil {
			log.Printf("Error flushing submit response: %v", err)
			return
		}
	}
}
