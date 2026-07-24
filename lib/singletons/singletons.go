package singletons

import "gqueue/v2/lib/execution"

var (
	TaskRegistry = execution.GetRegistry()
	TaskQueue    = execution.GetTaskQueue()
)
