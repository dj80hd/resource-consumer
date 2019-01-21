IMAGE_NAME=docker-automation.artifactory.uptake.com/automation/resource-consumer
COMMIT_SHA=$(shell git rev-parse HEAD)

default: build

.PHONY: dep test race cover docker migration

dep:
	dep ensure

format:
	go fmt ./...
	go vet `go list ./...`
	for pkg in `go list ./...`; do \
		golint -set_exit_status $$pkg || exit 1; \
	done

build: format
	go build -o consume-cpu/consume-cpu consume-cpu/consume_cpu.go
	go build -o consumer resource_consumer.go resource_consumer_handler.go utils.go

clean:
	rm -f consume-cpu/consume-cpu
	rm -f consumer

cover:
	goverage -covermode=set -coverprofile=cov.out `go list ./...`
	gocov convert cov.out | gocov report
	go tool cover -html=cov.out

docker: dep build
	docker build --rm --build-arg GIT_COMMIT="$(COMMIT_SHA)" --tag "$(IMAGE_NAME):latest" .

publish: docker
	docker push "$(IMAGE_NAME):latest"
