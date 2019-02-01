#! /usr/bin/env bash
set -euo pipefail

_report_usage() {
  echo "CPU: $(docker stats --no-stream | grep resource-consumer | awk '{print $3}')"
  echo "MEM: $(docker stats --no-stream | grep resource-consumer | awk '{print $4,$5,$6}' | tr -d ' ')"
  echo "DSK: $(docker ps -s | grep resource-consumer | awk '{print $(NF-2),$(NF-1),$NF}')"
}

_stop_container() {
  docker stop resource-consumer >/dev/null 2>&1 || true
  docker rm resource-consumer >/dev/null 2>&1 || true
}

_start_container() {
  docker run --name resource-consumer -d -p 8080:8080 dj80hd/resource-consumer
}

_set_load() {
curl --data "millicores=125&durationSec=600" http://localhost:8080/consume-cpu
curl --data "megabytes=200&durationSec=300" http://localhost:8080/consume-mem
curl --data "gigabytes=1&filename=/var/log/foo.log" http://localhost:8080/consume-disk
curl --data "metric=foo&delta=1.14&durationSec=300" http://localhost:8080/bump-metric
curl http://localhost:8080/bump-metric
}

_cleanup() {
  curl --data "gigabytes=0&filename=/var/log/foo.log" http://localhost:8080/consume-disk
}

_test() {
  _stop_container
  _start_container
  sleep 1
  _report_usage
  _set_load
  sleep 4
  _report_usage
  _cleanup
  _stop_container
  echo "Ok."
}

_load() {
  local host=${1:-localhost:8080}
  curl --data "millicores=250&durationSec=600" ${host}/consume-cpu
  curl --data "megabytes=500&durationSec=300" ${host}/consume-mem
  curl --data "gigabytes=4&filename=/var/log/rc-${RANDOM}.log" ${host}/consume-disk
  curl --data "metric=foo&delta=1.14&durationSec=300" ${host}/bump-metric
}

_deploy() {
  cat <<EOF | kubectl apply -f -
---
apiVersion: v1
kind: Service
metadata:
  labels:
    app: resource-consumer
  name: resource-consumer
spec:
  ports:
  - name: http
    port: 80
    protocol: TCP
    targetPort: 8080
  selector:
    app: resource-consumer
---
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: resource-consumer
  name: resource-consumer
spec:
  replicas: 1
  selector:
    matchLabels:
      app: resource-consumer
  template:
    metadata:
      name: resource-consumer
      labels:
        app: resource-consumer
    spec:
      containers:
      - image: dj80hd/resource-consumer
        imagePullPolicy: Always
        livenessProbe:
          httpGet:
            path: "/metrics"
            port: 8080
          initialDelaySeconds: 3
          periodSeconds: 10
        name: resource-consumer
        ports:
        - containerPort: 8080
          name: http
        resources:
          requests:
            cpu: 0.12
            memory: 200M
          limits:
            cpu: 2
            memory: 4Gi
---
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  annotations:
    uptake.run/repo: https://github.com/dj80hd/resource-consumer
  labels:
    app: resource-consumer
  name: resource-consumer
spec:
  rules:
  - host: resource-consumer.apps.infralab.mt2.uptake.run
    http:
      paths:
      - backend:
          serviceName: resource-consumer
          servicePort: 80
EOF
}
