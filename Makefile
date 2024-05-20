.PHONY=build

build-enqueuer:
	@go build -o bin/enqueuer cmd/enqueuer/enqueuer.go

build-process:
	@go build -o bin/process-job cmd/process_job/process_job.go

run-enqueuer: build-enqueuer
	@./bin/enqueuer

run-process: build-process
	@./bin/process-job

test:
	@go test -v -cover ./test/...