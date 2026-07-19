package network

import (
	"gqueue/v2/internal/execution"
	"log"
	"slices"

	"github.com/google/uuid"
)

func HandleSubmitTask(incomingMessage string, q *execution.TaskQueue, r *execution.TaskRegistry) {
	payloadParsed, err := FromCanonicalJson[TaskSubmission](incomingMessage)
	if err != nil {
		log.Printf("invalid submit message: %v", err)
		return
	}

	taskName := payloadParsed.TaskToExec.TaskName
	if !slices.Contains(r.List(), taskName) {
		log.Printf("task %s not in available tasks", taskName)
		return
	}

	if err := q.Submit(payloadParsed.TaskToExec.TaskName, payloadParsed.TaskToExec.Payload); err != nil {
		log.Printf("failed to submit task %s: %v", taskName, err)
	}
}

func HandlePullTask(incomingMessage string, q *execution.TaskQueue) string {
	payloadParsed, err := FromCanonicalJson[PullMesasges](incomingMessage)
	if err != nil {
		log.Printf("invalid pull message: %v", err)
		return ""
	}

	tasks := q.List(payloadParsed.Concurrency)
	envelopedData := make([]TaskEnvelope, 0, len(tasks))

	for _, t := range tasks {
		eid, err := uuid.NewV7()
		if err != nil {
			log.Printf("failed generating task id: %v", err)
			continue
		}

		envelopedData = append(envelopedData, TaskEnvelope{
			Type: "assign",
			Eid:  eid.String(),
			Task: TaskExecution{
				Payload:  t.Input,
				TaskName: t.Task,
			},
		})
	}

	return ToCanonicalJson(MultiTaskEnvelope{Type: "assign", Tasks: envelopedData})
}

func runTask(task TaskEnvelope, r *execution.TaskRegistry) {
	taskID := task.Eid
	taskHandler := task.Task.TaskName
	taskPayload := task.Task.Payload

	handler, err := r.GetHandler(taskHandler)
	if err != nil {
		log.Printf("Send NACK for taskId: %s reason: %v", taskID, err)
		return
	}

	runError := handler.Run(taskPayload)
	if runError != nil {
		log.Printf("Send NACK for taskId: %s reason: %v", taskID, runError)
		return
	}

	log.Printf("Send ACK for taskId: %s", taskID)
}

func HandleAssignTask(incomingMessage string, r *execution.TaskRegistry) {
	tasks, err := FromCanonicalJson[MultiTaskEnvelope](incomingMessage)
	if err != nil {
		log.Printf("invalid assign message: %v", err)
		return
	}
	for _, task := range tasks.Tasks {
		go runTask(task, r)
	}
}
