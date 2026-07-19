package network

import (
	"bufio"
	"log"
	"net"
)

func HandleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	msg, err := reader.ReadString('\n')
	if err != nil {
		log.Printf("Error: %s", err)
		return
	}
	log.Printf("Message received: %s", msg)
	msgStruct, err := FromCanonicalJson[Message](msg)
	if err != nil {
		log.Printf("Error: %s", err)
		return
	}
	switch msgStruct.Type {
	case "ack":
		log.Println("Got ack msg")
	default:
		log.Printf("Got msg type: %s", msgStruct.Type)
	}
}
