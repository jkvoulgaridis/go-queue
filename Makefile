build-example:
	go build  -o ./bin ./examples/cmd/...

submit-to-cluster:
	docker compose run --rm worker /app/bin/task_emitter -count 10000 -prob 10

get-metrics:
	curl localhost:8090/metrics