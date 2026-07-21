package nodes

import (
	"gqueue/v2/internal/network"
	"log"
	"net"
	"os"
)

func MasterSpinup() {
	addr := os.Getenv("BROKER_ADDR")
	if addr == "" {
		addr = ":18018"
	}

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("error listening on %s: %v", addr, err)
	}
	defer listener.Close()

	log.Printf("broker listening on %s", addr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("error accepting connection: %v", err)
			continue
		}
		go network.HandleConnection(conn)
	}
}
