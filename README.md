# Resource Consumer

A simplified version of [the kubernetes resource-consumer](https://github.com/kubernetes/kubernetes/tree/master/test/images/resource-consumer) that includes disk usage and is decoupled from the kubernetes build system.

## Overview
Resource Consumer is a tool which allows to generate cpu/memory/disk utilization in a container.
Resource Consumer can help with testing:
- cluster size autoscaling,
- horizontal autoscaling of pod - changing the size of replication controller,
- vertical autoscaling of pod - changing its resource limits.
- eviction scenarios

## Usage
Resource Consumer starts an HTTP server and handle sent requests.
It listens on port given as a flag (default 8080).
Action of consuming resources is send to the container by a POST http request.
Each http request creates new process.
Http request handler is in file resource_consumer_handler.go

The container consumes specified amount of resources:

- CPU in millicores
- Memory in megabytes
- Fake custom metrics
- Disk files in gigabytes

### Consume CPU http request
- suffix "ConsumeCPU"
- parameters "millicores" and "durationSec"

Consumes specified amount of millicores for durationSec seconds.
Consume CPU uses "./consume-cpu/consume-cpu" binary (file consume-cpu/consume_cpu.go).
When CPU consumption is too low this binary uses cpu by calculating math.sqrt(0) 10^7 times
and if consumption is too high binary sleeps for 10 millisecond.
One replica of Resource Consumer cannot consume more that 1 cpu.

### Consume Memory http request
- suffix "ConsumeMem"
- parameters "megabytes" and "durationSec"

Consumes specified amount of megabytes for durationSec seconds.
Consume Memory uses stress tool (stress -m 1 --vm-bytes megabytes --vm-hang 0 -t durationSec).
Request leading to consuming more memory then container limit will be ignored.

### Bump value of a fake custom metric
- suffix "BumpMetric"
- parameters "metric", "delta" and "durationSec"

Bumps metric with given name by delta for durationSec seconds.
Custom metrics in Prometheus format are exposed on "/metrics" endpoint.

### Consume disk request
- suffix "ConsumeDisk"
- parameters "gigabytes" and "filename"

Creates a filename whose size and name is specified by input.
Requests to create files in non-existent directories will be ignored.

### Running resource consumer
```bash
docker run -d -p 8080:8080 dj80hd/resource-consumer
```

### CURL examples

Take up 1/2 a CPU for 10 minutes:
```bash
curl --data "millicores=500&durationSec=600" http://localhost:8080/ConsumeCPU
```

Take up 1G Memory for 5 minutes:
```bash
curl --data "megabytes=1024&durationSec=300" http://localhost:8080/ConsumeMem
```

Take up 8G of disk:
```bash
curl --data "gigabytes=8&filename=/tmp/foo.txt" http://localhost:8080/ConsumeDisk
```

Free up same 8G of disk:
```bash
curl --data "gigabytes=0&filename=/tmp/foo.txt" http://localhost:8080/ConsumeDisk
```
