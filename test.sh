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
