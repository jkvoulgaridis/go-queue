package main

import (
	"gqueue/v2/examples"
	"log"
	"time"
)

func main() {
	if err := examples.ExampleFlow(); err != nil {
		log.Fatal(err)
	}

	time.Sleep(2 * time.Second)
}
