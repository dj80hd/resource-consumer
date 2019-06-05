VERSION=$(shell cat VERSION)
COMMIT_SHA=$(shell git rev-parse HEAD)
IMAGE_NAME=dj80hd/resource-consumer

default: build

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
	go build -o resource-consumer resource_consumer.go resource_consumer_handler.go utils.go

clean:
	rm -f consume-cpu/consume-cpu
	rm -f resource-consumer

docker: dep build
	docker build --rm --build-arg GIT_COMMIT="$(COMMIT_SHA)" --tag "$(IMAGE_NAME):$(VERSION)" .
	docker tag "$(IMAGE_NAME):$(VERSION)" "$(IMAGE_NAME):latest"

publish: docker
	docker push "$(IMAGE_NAME):$(VERSION)"
	docker push "$(IMAGE_NAME):latest"
