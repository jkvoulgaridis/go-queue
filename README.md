[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

# gqueue

`gqueue` is a lightweight distributed task queue library in Go.

🔴 Disclaimer: This is a toy project to get my hands dirty on Golang. 

It provides:
- A **broker** that accepts tasks and serves them to workers
- **Workers** that pull tasks and execute registered handlers
- A simple public API to register tasks and submit work over TCP
- Basic queue metrics (`queued`, `in_flight`, `failed`)

This project is designed for teams that want a simple, hackable queueing system without pulling in heavyweight infrastructure.

## Table of contents
- [Goals](#goals)
- [Architecture](#architecture)
- [How it works](#how-it-works)
- [Project structure](#project-structure)
- [Installation](#installation)
- [Quick start (local Docker)](#quick-start-local-docker)
- [Library usage](#library-usage)
- [Examples](#examples)
- [Configuration](#configuration)
- [Metrics](#metrics)
- [Development](#development)
- [Current behavior and limitations](#current-behavior-and-limitations)
- [Roadmap ideas](#roadmap-ideas)

## Goals
- Keep task queue setup and integration minimal.
- Expose a small Go API for producers and workers.
- Make behavior easy to reason about and easy to extend.
- Support basic horizontal worker scaling out of the box.

## Architecture

High-level components:

1. **Broker (`master`)**
   - Listens on TCP (default `:18018`)
   - Accepts submitted tasks
   - Stores queued / in-flight / failed state in memory
   - Serves task batches to workers on pull
   - Exposes HTTP metrics at `:8090/metrics`

2. **Worker**
   - Repeatedly pulls tasks from broker
   - Executes task handlers by name
   - Sends `ack` on success
   - Sends `nack` with reason on failure

3. **Producer / Emitter**
   - Submits tasks to broker via `SubmitTask(...)`

Communication protocol uses newline-delimited canonical JSON messages over TCP.

## How it works
1. Your app registers handlers (`taskName -> handler`).
2. Producers submit a task with a name + payload.
3. Broker appends task to in-memory queue.
4. Workers send pull requests with desired concurrency.
5. Broker returns tasks.
6. Worker runs the handler and responds with:
   - `ack` if successful (removed from `in_flight`)
   - `nack` if failed (moved to `failed`)

## Project structure
- `lib/public.go` — public API (`RegisterTasks`, `SubmitTask`, `RunBroker`, `RunWorker`)
- `lib/nodes/` — broker and worker runtime loops
- `lib/network/` — TCP handlers, protocol parsing, metrics endpoint adapter
- `lib/execution/` — task registry and queue internals
- `lib/messages/` — protocol message types
- `examples/` — sample handlers and runnable binaries (`broker`, `worker`, `task_emitter`)


### Build example binaries
```bash
make build-example
```

This produces:
- `bin/broker`
- `bin/worker`
- `bin/task_emitter`

## Quick start (local Docker)

Start broker + workers:
```bash
docker compose up --build
```

Submit example workload (from another terminal):
```bash
make submit-to-cluster
```

Read metrics:
```bash
make get-metrics
```

## Library usage

### 1) Define task handlers
A handler must implement:

```go
type Task interface {
    Run(string) error
}
```

### 2) Register tasks
```go
package main

import (
    "gqueue/v2/lib"
)

type EmailHandler struct{}

func (h EmailHandler) Run(payload string) error {
    // process payload
    return nil
}

func setup() error {
    return lib.RegisterTasks([]lib.TaskDef{
        {Name: "send_email", Handler: EmailHandler{}},
    })
}
```

### 3) Start broker
```go
func main() {
    if err := setup(); err != nil {
        panic(err)
    }
    lib.RunBroker()
}
```

### 4) Start worker
```go
func main() {
    if err := setup(); err != nil {
        panic(err)
    }
    lib.RunWorker()
}
```

### 5) Submit tasks from producer
```go
err := lib.SubmitTask("localhost:18018", "send_email", `{"to":"hello@example.com"}`)
if err != nil {
    panic(err)
}
```

## Examples

Runnable examples are provided in:
- `examples/cmd/broker/main.go`
- `examples/cmd/worker/main.go`
- `examples/cmd/task_emitter/main.go`
- `examples/handlers.go`

The included sample tasks are:
- `submit_report` (successful handler)
- `problematic` (always returns error)

This lets you quickly test success/failure behavior and observe queue metrics.

## Configuration

Environment variables:

- `BROKER_ADDR`
  - Broker bind address when running broker (default: `:18018`)
  - Broker endpoint when running workers/producers (example: `master:18018`)

- `WORKER_CONCURRENCY`
  - Number of tasks requested per worker pull cycle
  - Default: `1`

- `WORKER_POLL_MS`
  - Delay between worker pull cycles in milliseconds
  - Default: `2000`

## Metrics

Broker serves HTTP metrics at:
- `GET /metrics` on port `8090`

Current metric map keys:
- `queued`
- `in_flight`
- `failed`

Example:
```json
{"failed":3,"in_flight":20,"queued":107}
```

## Development

Useful commands:

```bash
make build-example
make submit-to-cluster
make get-metrics
```

Run binaries directly (without Docker):
```bash
BROKER_ADDR=:18018 ./bin/broker
BROKER_ADDR=localhost:18018 WORKER_CONCURRENCY=10 WORKER_POLL_MS=1000 ./bin/worker
BROKER_ADDR=localhost:18018 ./bin/task_emitter -count 100 -prob 10
```

## Current behavior and limitations
- Queue state is in-memory only (no persistence).
- No retry/backoff policy yet for failed tasks.
- No task timeout/dead-letter queue yet.
- Delivery semantics are simple pull + ack/nack (not exactly-once).
- Transport is raw TCP; auth/TLS are not implemented yet.

## Roadmap ideas
- Pluggable persistence (Redis/Postgres).
- Retry policies with exponential backoff.
- Delayed and scheduled tasks.
- Dead-letter queues.
- Graceful worker shutdown and task lease/visibility timeout.
- Optional TLS + authentication for broker endpoints.

---

If you use `gqueue` publicly, contributions and feedback are welcome.

## License

This project is licensed under the MIT License. See `LICENSE`.
