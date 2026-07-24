package execution

import (
	"errors"
	"fmt"
)

type Task interface {
	Run(string) error
}

type TaskRegistry struct {
	taskMap map[string]Task
}

func NewRegistry() *TaskRegistry {
	return &TaskRegistry{
		taskMap: make(map[string]Task),
	}
}

func (r *TaskRegistry) Register(name string, handler Task) error {
	if name == "" {
		return errors.New("task name cannot be empty")
	}

	if _, exists := r.taskMap[name]; exists {
		return fmt.Errorf("task %q is already registered", name)
	}

	r.taskMap[name] = handler
	return nil
}

func (r *TaskRegistry) List() []string {
	tasks := make([]string, 0, len(r.taskMap))
	for name := range r.taskMap {
		tasks = append(tasks, name)
	}
	return tasks
}

func (r *TaskRegistry) GetHandler(taskname string) (Task, error) {
	handler, exists := r.taskMap[taskname]
	if !exists {
		return nil, fmt.Errorf("task %q is not registered", taskname)
	}
	return handler, nil
}