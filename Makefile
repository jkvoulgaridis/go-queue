build-example:
	go build  -o ./bin ./examples/cmd/...

submit-to-cluster:
	docker compose run --rm worker /app/bin/task_emitter -count 100 -prob 0

get-metrics:
	curl localhost:8090/metrics