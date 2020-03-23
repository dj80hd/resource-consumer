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
COPY bin/stress /stress
RUN chmod 755 /stress

COPY --from=build /consume-cpu /consume-cpu

COPY --from=build /consumer /consumer
RUN  echo "/consumer" > /run.sh && chmod +x /run.sh
CMD /run.sh

EXPOSE 8080
