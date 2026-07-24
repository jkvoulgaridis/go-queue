package messages

type Message struct {
	Type string `json:"type"`
}

type PullMesasges struct {
	Type        string `json:"type"`
	Concurrency int    `json:"concurrency"`
}

type TaskExecution struct {
	Payload  string `json:"payload"`
	TaskName string `json:"task"`
}

type TaskSubmission struct {
	Type       string        `json:"type"`
	TaskToExec TaskExecution `json:"task"`
}

type TaskEnvelope struct {
	Type string        `json:"type"`
	Eid  string        `json:"id"`
	Task TaskExecution `json:"task"`
}

type MultiTaskEnvelope struct {
	Type  string         `json:"type"`
	Tasks []TaskEnvelope `json:"tasks"`
}

type AckMessage struct {
	Type string `json:"type"`
	Eid  string `json:"id"`
}

type NackMessage struct {
	Type   string `json:"type"`
	Eid    string `json:"id"`
	Reason string `json:"reason"`
}
