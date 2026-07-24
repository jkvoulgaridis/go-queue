package execution

import (
	"errors"
	"gqueue/v2/lib/messages"
	"log"
	"slices"
	"sync"

	"github.com/google/uuid"
)

type TaskQueueItem struct {
	Task  string
	Input string
}

type TaskQueue struct {
	q_lock      sync.Mutex
	i_lock      sync.Mutex
	f_lock      sync.Mutex
	queued      []messages.TaskEnvelope
	inFlight    []uuid.UUID
	failedTasks []uuid.UUID
}

func newTaskQueue() *TaskQueue {
	return &TaskQueue{
		queued:   make([]messages.TaskEnvelope, 0),
		inFlight: make([]uuid.UUID, 0),
	}
}

var (
	taskQ *TaskQueue
	onceT sync.Once
)

func GetTaskQueue() *TaskQueue {
	onceT.Do(func() {
		taskQ = newTaskQueue()
	})
	return taskQ
}

func (r *TaskQueue) addTaskEnvelope(task TaskQueueItem) messages.TaskEnvelope {
	eid, err := uuid.NewV7()
	if err != nil {
		log.Printf("failed generating task id: %v", err)
		panic(err)
	}
	enveloped := messages.TaskEnvelope{
		Type: "assign",
		Eid:  eid.String(),
		Task: messages.TaskExecution{Payload: task.Input, TaskName: task.Task},
	}
	return enveloped
}

func (r *TaskQueue) Submit(name string, payload string) error {
	if name == "" {
		return errors.New("task name cannot be empty")
	}

	task := TaskQueueItem{
		Task:  name,
		Input: payload,
	}
	enveloped := r.addTaskEnvelope(task)
	r.q_lock.Lock()
	r.queued = append(r.queued, enveloped)
	r.q_lock.Unlock()
	return nil
}

func (r *TaskQueue) Execute() messages.TaskEnvelope {
	if len(r.queued) == 0 {
		return messages.TaskEnvelope{}
	}
	r.q_lock.Lock()
	item := r.queued[0]
	r.queued[0] = messages.TaskEnvelope{}
	r.queued = r.queued[1:]
	r.q_lock.Unlock()

	return item
}

func (r *TaskQueue) List(n int) []messages.TaskEnvelope {
	if n <= 0 || len(r.queued) == 0 {
		return []messages.TaskEnvelope{}
	}
	if n > len(r.queued) {
		n = len(r.queued)
	}
	items := r.queued[:n]
	r.queued = r.queued[n:]
	return items
}

func (r *TaskQueue) AddToInFlight(i uuid.UUID) {
	r.i_lock.Lock()
	r.inFlight = append(r.inFlight, i)
	r.i_lock.Unlock()
}

func (r *TaskQueue) GetMetrics() map[string]int {
	met := map[string]int{
		"queued":    len(r.queued),
		"in_flight": len(r.inFlight),
		"failed":    len(r.failedTasks),
	}
	return met
}

func (r *TaskQueue) RemoveFromInFlight(eid uuid.UUID, isFailed bool) {

	idx := slices.Index(r.inFlight, eid)
	if idx != -1 {
		r.i_lock.Lock()
		r.inFlight = slices.Delete(r.inFlight, idx, idx+1)
		r.i_lock.Unlock()
		log.Printf("Removed item: %s\n", eid)
	}
	if isFailed {
		r.f_lock.Lock()
		r.failedTasks = append(r.failedTasks, eid)
		r.f_lock.Unlock()
	}
}
