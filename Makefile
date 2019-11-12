VERSION=$(shell cat VERSION)
COMMIT_SHA=$(shell git rev-parse HEAD)
IMAGE_NAME=myspot
LINTER_INSTALLED := $(shell sh -c 'which golangci-lint')

default: build

lint:
	go fmt ./...
	go vet `go list ./...`
ifdef LINTER_INSTALLED
	golangci-lint run
else
	$(error golangci-lint not found, skipping linting. Installation instructions: https://github.com/golangci/golangci-lint#ci-installation)
endif

build: lint
	go build -o consume-cpu/consume-cpu consume-cpu/consume_cpu.go
	go build -o resource-consumer resource_consumer.go resource_consumer_handler.go utils.go

clean:
	rm -f consume-cpu/consume-cpu
	rm -f resource-consumer

docker: build
	docker build --rm --build-arg GIT_COMMIT="$(COMMIT_SHA)" --tag "$(IMAGE_NAME):$(VERSION)" .
	docker tag "$(IMAGE_NAME):$(VERSION)" "$(IMAGE_NAME):latest"

publish: docker
	docker push "$(IMAGE_NAME):$(VERSION)"
	docker push "$(IMAGE_NAME):latest"
