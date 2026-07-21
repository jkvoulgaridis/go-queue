package main

import (
	"flag"
	"gqueue/v2/examples"
)

func main() {
	mode := flag.String("mode", "", "run mode: broker|worker|submit")
	brokerAddr := flag.String("broker", "", "broker address for submit mode (default BROKER_ADDR or localhost:18018)")
	submitCount := flag.Int("count", 2, "number of tasks to submit in submit mode")
	problematicPercent := flag.Float64("problematic-percent", 50, "percentage (0-100) of submitted tasks that should be problematic")
	flag.Parse()

	examples.App(*mode, *brokerAddr, *submitCount, *problematicPercent)
}
