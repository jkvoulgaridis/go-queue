package execution

import (
	"errors"
)

type TaskQueueItem struct {
	Task  string
	Input string
}

type TaskQueue struct {
	taskQ []TaskQueueItem
}

func NewTaskQueue() *TaskQueue {
	return &TaskQueue{
		taskQ: make([]TaskQueueItem, 0),
	}
}

func (r *TaskQueue) Submit(name string, payload string) error {
	if name == "" {
		return errors.New("task name cannot be empty")
	}

	r.taskQ = append(r.taskQ, TaskQueueItem{
		Task:  name,
		Input: payload,
	})
	return nil
}

func (r *TaskQueue) Execute() TaskQueueItem {
	if len(r.taskQ) == 0 {
		return TaskQueueItem{}
	}
	item := r.taskQ[0]
	r.taskQ[0] = TaskQueueItem{}
	r.taskQ = r.taskQ[1:]
	return item
}

func (r *TaskQueue) List(n int) []TaskQueueItem {
	if n <= 0 || len(r.taskQ) == 0 {
		return []TaskQueueItem{}
	}
	if n > len(r.taskQ) {
		n = len(r.taskQ)
	}
	items := r.taskQ[:n]
	r.taskQ = r.taskQ[n:]
	return items
}
