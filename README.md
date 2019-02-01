# Resource Consumer

A simplified version of [the kubernetes resource-consumer](https://github.com/kubernetes/kubernetes/tree/master/test/images/resource-consumer) that includes disk usage and is decoupled from the kubernetes build system.

## Overview
Resouce Consumer allows one to generate the following type of load inside a container:
- CPU in millicores
- Memory in megabytes
- Fake custom metrics
- Disk files in gigabytes

Resource Consumer can help with testing:
- cluster size autoscaling,
- horizontal autoscaling of pod - changing the size of replication controller,
- vertical autoscaling of pod - changing its resource limits.
- eviction scenarios
- reserved resources

## Usage

Run Resource Consumer locally:
```bash
docker run --name resource-consumer -d -p 8080:8080 dj80hd/resource-consumer
```

## CURL examples

* Take up 1/8 CPU for 10 minutes:
```bash
curl --data "millicores=125&durationSec=600" http://localhost:8080/consume-cpu
```
Note: One replica of Resource Consumer cannot consume more that 1 cpu.

* Take up 200M Memory for 5 minutes:
```bash
curl --data "megabytes=200&durationSec=300" http://localhost:8080/consume-mem
```
Note: Request to consume more memory then container limit will be ignored.

* Take up 1G of disk:
```bash
curl --data "gigabytes=1&filename=/var/log/foo.log" http://localhost:8080/consume-disk
```
Note: Requests to create files in non-existent directories will be ignored.

* Free up same 1G of disk:
```bash
curl --data "gigabytes=0&filename=/var/log/foo.log" http://localhost:8080/consume-disk
```
Note: Size of 0 gigabytes deletes file

* Set metric `foo` to 1.14 for 5 minutes:
```bash
curl --data "metric=foo&delta=1.14&durationSec=300" http://localhost:8080/bump-metric
```
Note: Custom metrics in Prometheus format are exposed on "/metrics" endpoint.

## Test

Observe local cpu, mem, and disk:
```bash
echo "CPU: $(docker stats --no-stream | grep resource-consumer | awk '{print $3}')" && \
echo "MEM: $(docker stats --no-stream | grep resource-consumer | awk '{print $4,$5,$6}' | tr -d ' ')" && \
echo "DSK: $(docker ps -s | grep resource-consumer | awk '{print $(NF-2),$(NF-1),$NF}')"
```
