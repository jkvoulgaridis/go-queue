package nodes

import (
	"fmt"
	"gqueue/v2/lib/network"
	"log"
	"net"
	"net/http"
	"os"
)

func SpinupTCPChannel() {
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

func getMetrics(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "OK\n")
}

func SpinUpHTTPMetrics() {
	http.HandleFunc("/metrics", getMetrics)
	http.ListenAndServe(":8090", nil)

}

func MasterSpinup() {
	SpinupTCPChannel()
	go SpinUpHTTPMetrics()
}
