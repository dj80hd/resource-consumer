# build stage
FROM golang:1.10-alpine AS build

# golang base image has GOPATH=/go
ADD . /go/src/github.com/dj80hd/resource-consumer
WORKDIR /go/src/github.com/dj80hd/resource-consumer

RUN GOOS=linux GOARCH=amd64 go build -o /consume-cpu consume-cpu/consume_cpu.go
RUN GOOS=linux GOARCH=amd64 go build -o /consumer resource_consumer.go resource_consumer_handler.go utils.go


ENV STRESS_VERSION=1.0.4 
RUN \
  apk add --update bash g++ make curl && \
  curl -o /tmp/stress-${STRESS_VERSION}.tgz https://people.seas.harvard.edu/~apw/stress/stress-${STRESS_VERSION}.tar.gz && \
  cd /tmp && tar xvf stress-${STRESS_VERSION}.tgz && rm /tmp/stress-${STRESS_VERSION}.tgz && \
  cd /tmp/stress-${STRESS_VERSION} && \
  ./configure && make && make install && \
  apk del g++ make curl && \
  rm -rf /tmp/* /var/tmp/* /var/cache/apk/* /var/cache/distfiles/*

# image stage
FROM alpine:latest

# stress tool
ENV STRESS_VERSION=1.0.4
RUN \
  apk add --update bash g++ make curl && \
  curl -o /tmp/stress-${STRESS_VERSION}.tgz https://people.seas.harvard.edu/~apw/stress/stress-${STRESS_VERSION}.tar.gz && \
  cd /tmp && tar xvf stress-${STRESS_VERSION}.tgz && rm /tmp/stress-${STRESS_VERSION}.tgz && \
  cd /tmp/stress-${STRESS_VERSION} && \
  ./configure && make && make install && \
  apk del g++ make curl && \
  rm -rf /tmp/* /var/tmp/* /var/cache/apk/* /var/cache/distfiles/*

COPY --from=build /consumer /consumer
COPY --from=build /consume-cpu /consume-cpu
EXPOSE 8080
CMD ["/consumer"]
