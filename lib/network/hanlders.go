package network

import (
	"gqueue/v2/lib/messages"
	"gqueue/v2/lib/singletons"
	"log"
	"slices"

	"github.com/google/uuid"
)

func HandleSubmitTask(incomingMessage string) {
	q := singletons.TaskQueue
	r := singletons.TaskRegistry
	payloadParsed, err := FromCanonicalJson[messages.TaskSubmission](incomingMessage)
	if err != nil {
		log.Printf("broker submit parse_error err=%v raw=%q", err, incomingMessage)
		return
	}

	taskName := payloadParsed.TaskToExec.TaskName
	payload := payloadParsed.TaskToExec.Payload
	if !slices.Contains(r.List(), taskName) {
		log.Printf("broker submit rejected task=%s reason=unregistered payload=%q", taskName, payload)
		return
	}

	if err := q.Submit(taskName, payload); err != nil {
		log.Printf("broker submit enqueue_error task=%s payload=%q err=%v", taskName, payload, err)
		return
	}

	metrics := q.GetMetrics()
	log.Printf("broker submit queued task=%s payload=%q queued=%d in_flight=%d", taskName, payload, metrics["queued"], metrics["in_flight"])
}

func HandlePullTask(incomingMessage string) string {
	q := singletons.TaskQueue
	payloadParsed, err := FromCanonicalJson[messages.PullMesasges](incomingMessage)
	if err != nil {
		log.Printf("broker pull parse_error err=%v raw=%q", err, incomingMessage)
		return ""
	}

	if payloadParsed.Concurrency <= 0 {
		log.Printf("broker pull invalid_concurrency=%d", payloadParsed.Concurrency)
		payloadParsed.Concurrency = 1
	}

	tasks := q.List(payloadParsed.Concurrency)
	envelopedData := make([]messages.TaskEnvelope, 0, len(tasks))

	for _, t := range tasks {
		eid, _ := uuid.Parse(t.Eid)
		q.AddToInFlight(eid)
		envelopedData = append(envelopedData, t)
	}

	metrics := q.GetMetrics()
	log.Printf("broker pull assigned requested=%d assigned=%d queued=%d in_flight=%d", payloadParsed.Concurrency, len(envelopedData), metrics["queued"], metrics["in_flight"])

	return ToCanonicalJson(messages.MultiTaskEnvelope{Type: "assign", Tasks: envelopedData})
}

func HandleAckMessage(incomingMessage messages.AckMessage) {
	q := singletons.TaskQueue
	taskId, _ := uuid.Parse(incomingMessage.Eid)
	if taskId == uuid.Nil {
		log.Printf("broker ack invalid_id raw=%q", incomingMessage.Eid)
		return
	}
	q.RemoveFromInFlight(taskId, false)
	metrics := q.GetMetrics()
	log.Printf("broker ack task_id=%s queued=%d in_flight=%d", incomingMessage.Eid, metrics["queued"], metrics["in_flight"])
}

func HandleNackkMessage(incomingMessage messages.NackMessage) {
	q := singletons.TaskQueue
	taskId, _ := uuid.Parse(incomingMessage.Eid)
	if taskId == uuid.Nil {
		log.Printf("broker ack invalid_id raw=%q", incomingMessage.Eid)
		return
	}
	q.RemoveFromInFlight(taskId, true)
	metrics := q.GetMetrics()
	log.Printf("broker ack task_id=%s queued=%d in_flight=%d", incomingMessage.Eid, metrics["queued"], metrics["in_flight"])
}
