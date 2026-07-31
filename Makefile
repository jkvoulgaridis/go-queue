TEST_PKGS := $(shell go list ./... | grep -v '/examples')

build-example:
	go build  -o ./bin ./examples/cmd/...

submit-to-cluster:
	docker compose run --rm worker /app/bin/task_emitter -count 10000 -prob 10

get-metrics:
	curl localhost:8090/metrics

test:
	go test -cover $(TEST_PKGS)

test-cov-out:
	go test -coverprofile=coverage.out $(TEST_PKGS)
	go tool cover -html=coverage.out

total-coverage:
	go tool cover -func=coverage.out | awk '/^total:/ {print $$3}'
