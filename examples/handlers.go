package examples

import (
	"errors"
	"log"
	"time"
)

type SubmitReportHandler struct{}

type ProblematicHandler struct{}

func (h SubmitReportHandler) Run(reportPayload string) error {
	log.Printf("Submitting report: %s", reportPayload)
	time.Sleep(200 * time.Millisecond)
	log.Printf("Submitted report: %s", reportPayload)
	return nil
}

func (h ProblematicHandler) Run(_ string) error {
	time.Sleep(200 * time.Millisecond)
	return errors.New("random issue")
}
