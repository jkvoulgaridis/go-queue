package main

import (
	"gqueue/v2/examples"
	"gqueue/v2/lib"
)

func main() {
	examples.SetupTasks()
	lib.RunWorker()
}
