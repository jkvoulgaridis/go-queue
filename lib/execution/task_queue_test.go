package execution

import (
	"testing"

	"github.com/google/uuid"
)

func TestExecuteFromEmptyQueue(t *testing.T) {
	q := GetTaskQueue()
	q.Purge()

	item := q.Execute()

	if item.Eid != "" || item.Type != "" || item.Task.TaskName != "" || item.Task.Payload != "" {
		t.Errorf("expected empty envelope for empty queue, got: %#v", item)
	}
}

func TestExecuteDequeuesOneItem(t *testing.T) {
	q := GetTaskQueue()
	q.Purge()

	err := q.Submit("exec-task", "payload")
	if err != nil {
		t.Fatalf("unexpected submit error: %v", err)
	}

	item := q.Execute()

	if item.Task.TaskName != "exec-task" || item.Task.Payload != "payload" {
		t.Errorf("unexpected executed item: %#v", item)
	}
	if len(q.queued) != 0 {
		t.Errorf("expected queue to be empty after execute, got: %d", len(q.queued))
	}
}

func TestTaskQueueCreate(t *testing.T) {
	q := GetTaskQueue()

	if len(q.failedTasks) > 0 {
		t.Errorf("expected failed tasks to be 0 on start")
	}

	if len(q.inFlight) > 0 {
		t.Errorf("expected in flight tasks to be 0 on start")
	}

	if len(q.queued) > 0 {
		t.Errorf("expected queued tasks to be 0 on start")
	}
}

func TestSubmitEmptyTask(t *testing.T) {
	q := GetTaskQueue()
	q.Purge()

	err := q.Submit("", "")
	if err == nil {
		t.Errorf("expected error for empty task name")
	}
	if len(q.queued) > 0 {
		t.Errorf("expected queued tasks to be 0 for empty task")
	}
}

func TestSubmitTask(t *testing.T) {
	q := GetTaskQueue()
	q.Purge()

	err := q.Submit("test-task", "")
	if err != nil {
		t.Fatalf("unexpected submit error: %v", err)
	}
	if len(q.queued) != 1 {
		t.Errorf("expected queued tasks to be 1")
	}
	task := TaskQueueItem{
		Task:  "test-task",
		Input: "",
	}
	taskEnvelope := q.addTaskEnvelope(
		task,
	)
	if task.Input != taskEnvelope.Task.Payload || task.Task != taskEnvelope.Task.TaskName {
		t.Errorf("Expected task be have been packed in envelope immutable")
	}
}

func TestListEdgeCases(t *testing.T) {
	q := GetTaskQueue()
	q.Purge()

	if err := q.Submit("t1", "p1"); err != nil {
		t.Fatalf("unexpected submit error: %v", err)
	}
	if err := q.Submit("t2", "p2"); err != nil {
		t.Fatalf("unexpected submit error: %v", err)
	}

	if len(q.List(0)) != 0 {
		t.Errorf("expected empty list for n=0")
	}
	if len(q.List(-1)) != 0 {
		t.Errorf("expected empty list for n<0")
	}

	items := q.List(10)
	if len(items) != 2 {
		t.Errorf("expected to receive all queued items, got: %d", len(items))
	}
	if len(q.queued) != 0 {
		t.Errorf("expected queue drained, got: %d", len(q.queued))
	}
}

func TestRemoveFromInFlightUnknownFailedTracksFailure(t *testing.T) {
	q := GetTaskQueue()
	q.Purge()

	unknown := uuid.New()
	q.RemoveFromInFlight(unknown, true)

	if len(q.inFlight) != 0 {
		t.Errorf("expected inFlight unchanged, got: %d", len(q.inFlight))
	}
	if len(q.failedTasks) != 1 {
		t.Errorf("expected failedTasks to record failed id, got: %d", len(q.failedTasks))
	}
	if q.failedTasks[0] != unknown {
		t.Errorf("expected failed id to match unknown id")
	}
}

func TestPurgeClearsAllCollections(t *testing.T) {
	q := GetTaskQueue()
	q.Purge()

	if err := q.Submit("cleanup-task", ""); err != nil {
		t.Fatalf("unexpected submit error: %v", err)
	}
	item := q.Execute()
	taskID, err := uuid.Parse(item.Eid)
	if err != nil {
		t.Fatalf("failed to parse task id: %v", err)
	}
	q.AddToInFlight(taskID)
	q.RemoveFromInFlight(taskID, true)

	q.Purge()

	if len(q.queued) != 0 || len(q.inFlight) != 0 || len(q.failedTasks) != 0 {
		t.Errorf("expected purge to clear all collections")
	}
}

func TestInflighHandling(t *testing.T) {
	q := GetTaskQueue()

	q.Purge()
	q.Submit("test-task", "")

	task := q.List(10)

	if len(task) != 1 {
		t.Errorf("Expected 1 task in the queue, received: %d", len(task))
	}
	taskId, _ := uuid.Parse(task[0].Eid)
	q.AddToInFlight(taskId)

	if len(q.inFlight) != 1 && len(q.queued) != 0 {
		t.Errorf("Expected different task distribution")
	}

	q.RemoveFromInFlight(taskId, false)

	if len(q.inFlight) != 0 && len(q.queued) != 0 && len(q.failedTasks) != 0 {
		t.Errorf("Expected different task distribution")
	}

	q.AddToInFlight(taskId)

	q.RemoveFromInFlight(taskId, true)
	if len(q.inFlight) != 0 && len(q.queued) != 0 && len(q.failedTasks) != 1 {
		t.Errorf("Expected different task distribution")
	}
}

func TestMetrics(t *testing.T) {
	q := GetTaskQueue()

	q.Purge()
	q.Submit("test-task", "")
	q.Submit("test-2", "")
	q.Submit("test-task-3", "")

	task := q.List(2)
	taskId1, _ := uuid.Parse(task[0].Eid)
	taskId2, _ := uuid.Parse(task[1].Eid)

	q.AddToInFlight(taskId1)
	q.AddToInFlight(taskId2)

	q.RemoveFromInFlight(taskId2, true)
	metrics := q.GetMetrics()

	if metrics["queued"] != 1 || metrics["in_flight"] != 1 || metrics["failed"] != 1 {
		t.Errorf("Got invalid metrics : %#v", metrics)
	}
}
