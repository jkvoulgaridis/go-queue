package execution

import (
	"errors"
	"log"
	"slices"

	"github.com/google/uuid"
)

type TaskQueueItem struct {
	Task  string
	Input string
}

type TaskQueue struct {
	queued   []TaskQueueItem
	inFlight []uuid.UUID
}

func NewTaskQueue() *TaskQueue {
	return &TaskQueue{
		queued:   make([]TaskQueueItem, 0),
		inFlight: make([]uuid.UUID, 0),
	}
}

func (r *TaskQueue) Submit(name string, payload string) error {
	if name == "" {
		return errors.New("task name cannot be empty")
	}

	r.queued = append(r.queued, TaskQueueItem{
		Task:  name,
		Input: payload,
	})
	return nil
}

func (r *TaskQueue) Execute() TaskQueueItem {
	if len(r.queued) == 0 {
		return TaskQueueItem{}
	}
	item := r.queued[0]
	r.queued[0] = TaskQueueItem{}
	r.queued = r.queued[1:]
	return item
}

func (r *TaskQueue) List(n int) []TaskQueueItem {
	if n <= 0 || len(r.queued) == 0 {
		return []TaskQueueItem{}
	}
	if n > len(r.queued) {
		n = len(r.queued)
	}
	items := r.queued[:n]
	r.queued = r.queued[n:]
	return items
}

func (r *TaskQueue) AddToInFlight(i uuid.UUID) {
	r.inFlight = append(r.inFlight, i)
}

func (r *TaskQueue) GetMetrics() map[string]int {
	met := map[string]int{
		"queued":    len(r.queued),
		"in_flight": len(r.inFlight),
	}
	return met
}

func (r *TaskQueue) RemoveFromInFlight(eid uuid.UUID) {
	idx := slices.Index(r.inFlight, eid)
	if idx != -1 {
		r.inFlight = slices.Delete(r.inFlight, idx, idx+1)
		log.Printf("Removed item: %s\n", eid) // [submit_report cleanup]
	}
}
