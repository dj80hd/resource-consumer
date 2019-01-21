# build stage
FROM golang:1.10-alpine AS build

# golang base image has GOPATH=/go
ADD . /go/src/github.com/dj80hd/resource-consumer

RUN CGO_ENABLED=0 go install github.com/dj80hd/resource-consumer/...

ENV STRESS_VERSION=1.0.4 
RUN \
  apk add --update bash g++ make curl && \
  curl -o /tmp/stress-${STRESS_VERSION}.tgz https://people.seas.harvard.edu/~apw/stress/stress-${STRESS_VERSION}.tar.gz && \
  cd /tmp && tar xvf stress-${STRESS_VERSION}.tgz && rm /tmp/stress-${STRESS_VERSION}.tgz && \
  cd /tmp/stress-${STRESS_VERSION} && \
  ./configure && make && make install && \
  apk del g++ make curl && \
  rm -rf /tmp/* /var/tmp/* /var/cache/apk/* /var/cache/distfiles/*

# actual image stage
FROM docker.artifactory.uptake.com/base/alpine:latest

#stress
ENV STRESS_VERSION=1.0.4 
RUN \
  apk add --update bash g++ make curl && \
  curl -o /tmp/stress-${STRESS_VERSION}.tgz https://people.seas.harvard.edu/~apw/stress/stress-${STRESS_VERSION}.tar.gz && \
  cd /tmp && tar xvf stress-${STRESS_VERSION}.tgz && rm /tmp/stress-${STRESS_VERSION}.tgz && \
  cd /tmp/stress-${STRESS_VERSION} && \
  ./configure && make && make install && \
  apk del g++ make curl && \
  rm -rf /tmp/* /var/tmp/* /var/cache/apk/* /var/cache/distfiles/*

ADD consumer /consumer
ADD consume-cpu /consume-cpu
ADD me-cpu /consume-cpu
EXPOSE 8080
ENTRYPOINT ["/consumer"]
