package main

import (
	"flag"
	"fmt"
	"gqueue/v2/lib"
	"log"
	"math"
	"os"
)

func SubmitSampleTasks(
	submitCount int,
	problematicPercent float64,
) {
	addr := os.Getenv("BROKER_ADDR")
	if submitCount < 0 {
		log.Fatal("count must be >= 0")
	}
	if problematicPercent < 0 || problematicPercent > 100 {
		log.Fatal("problematic-percent must be between 0 and 100")
	}

	problematicCount := int(math.Round(float64(submitCount) * (problematicPercent / 100.0)))
	if problematicCount > submitCount {
		problematicCount = submitCount
	}
	normalCount := submitCount - problematicCount

	for i := 0; i < normalCount; i++ {
		payload := fmt.Sprintf("q2-finance-report-%d", i+1)
		if err := lib.SubmitTask(addr, "submit_report", payload); err != nil {
			log.Fatal(err)
		}
	}
	for i := 0; i < problematicCount; i++ {
		payload := fmt.Sprintf("payload-%d", i+1)
		if err := lib.SubmitTask(addr, "problematic", payload); err != nil {
			log.Fatal(err)
		}
	}

	log.Printf("submitted tasks to broker=%s total=%d normal=%d problematic=%d problematic_percent=%.2f", addr, submitCount, normalCount, problematicCount, problematicPercent)
}

func main() {
	cnt := flag.Int("count", 100, "Number of tasks to fire")
	prct := flag.Float64("prob", 0.0, "Problemetic tasks to fire")
	flag.Parse()
	SubmitSampleTasks(*cnt, *prct)
}
