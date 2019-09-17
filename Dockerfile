# build stage
FROM golang:1.11-alpine AS build

ENV STRESS_VERSION=1.0.4
RUN apk add --update git g++ make

ADD . /go/src/github.com/dj80hd/resource-consumer
WORKDIR /go/src/github.com/dj80hd/resource-consumer
RUN GO111MODULE=on GOOS=linux GOARCH=amd64 go build -o /consume-cpu consume-cpu/consume_cpu.go
RUN GO111MODULE=on GOOS=linux GOARCH=amd64 go build -o /consumer resource_consumer.go resource_consumer_handler.go utils.go
# image stage
FROM alpine:latest

# stress tool
COPY bin/stress /usr/local/bin/stress
RUN chmod 644 /usr/local/bin/stress
COPY --from=build /consumer /consumer
COPY --from=build /consume-cpu /consume-cpu
EXPOSE 8080
CMD ["/consumer"]
