/*
Copyright 2015 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"
)

//ResourceConsumerHandler holds all state for the Handler
type ResourceConsumerHandler struct {
	metrics     map[string]float64
	metricsLock sync.Mutex
}

//NewResourceConsumerHandler is a constructor
func NewResourceConsumerHandler() *ResourceConsumerHandler {
	return &ResourceConsumerHandler{metrics: map[string]float64{}}
}

func (handler *ResourceConsumerHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// handle exposing metrics in Prometheus format (both GET & POST)
	if req.URL.Path == "/metrics" {
		handler.handleMetrics(w)
		return
	}
	if req.Method != "POST" {
		http.Error(w, "HTTP Post required", http.StatusBadRequest)
		return
	}

	if err := req.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.URL.Path == "/consume-cpu" {
		handler.handleConsumeCPU(w, req.Form)
		return
	}
	if req.URL.Path == "/consume-mem" {
		handler.handleConsumeMem(w, req.Form)
		return
	}
	if req.URL.Path == "/consume-disk" {
		handler.handleConsumeDisk(w, req.Form)
		return
	}
	if req.URL.Path == "/bump-metric" {
		handler.handleBumpMetric(w, req.Form)
		return
	}
	http.Error(w, fmt.Sprintf("%s: %s", "unknown", req.URL.Path), http.StatusNotFound)
}

func (handler *ResourceConsumerHandler) handleConsumeCPU(w http.ResponseWriter, query url.Values) {
	durationSecString := query.Get("durationSec")
	millicoresString := query.Get("millicores")

	if durationSecString == "" || millicoresString == "" {
		http.Error(w, "not", http.StatusBadRequest)
		return
	}

	durationSec, durationSecError := strconv.Atoi(durationSecString)
	millicores, millicoresError := strconv.Atoi(millicoresString)

	if durationSecError != nil || millicoresError != nil {
		http.Error(w, "non-integer millicores or durationSec", http.StatusBadRequest)
		return
	}

	go ConsumeCPU(millicores, durationSec)
	_, err := fmt.Fprintf(w, "/consume-cpu\nmillicores %d\ndurationSec %d\n", millicores, durationSec)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func (handler *ResourceConsumerHandler) handleConsumeDisk(w http.ResponseWriter, query url.Values) {
	gigabytesString := query.Get("gigabytes")
	durationSecString := query.Get("durationSec")
	filename := query.Get("filename")

	if gigabytesString == "" || durationSecString == "" || filename == "" {
		http.Error(w, "gigabytes, durationSec, or filename missing", http.StatusBadRequest)
		return
	}

	gigabytes, gigabytesError := strconv.Atoi(gigabytesString)
	if gigabytesError != nil {
		http.Error(w, "incorrect gigabytes", http.StatusBadRequest)
		return
	}

	durationSec, durationSecError := strconv.Atoi(durationSecString)
	if durationSecError != nil {
		http.Error(w, "durantionSec be integer", http.StatusBadRequest)
		return
	}

	go ConsumeDisk(gigabytes, durationSec, filename)
	_, err := fmt.Fprintf(w, "/consume-disk\ngigabytes %d\ndurationSec %d\nfilename %s\n", gigabytes, durationSec, filename)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func (handler *ResourceConsumerHandler) handleConsumeMem(w http.ResponseWriter, query url.Values) {
	durationSecString := query.Get("durationSec")
	megabytesString := query.Get("megabytes")

	if durationSecString == "" || megabytesString == "" {
		http.Error(w, "durantionSec or megabytes missing", http.StatusBadRequest)
		return
	}

	durationSec, durationSecError := strconv.Atoi(durationSecString)
	megabytes, megabytesError := strconv.Atoi(megabytesString)

	if durationSecError != nil || megabytesError != nil {
		http.Error(w, "durantionSec and megabytes must be integers", http.StatusBadRequest)
		return
	}

	go ConsumeMem(megabytes, durationSec)
	_, err := fmt.Fprintf(w, "/consume-mem\nmegabytes %d\ndurationSec %d\n", megabytes, durationSec)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func (handler *ResourceConsumerHandler) handleMetrics(w http.ResponseWriter) {
	handler.metricsLock.Lock()
	defer handler.metricsLock.Unlock()
	for k, v := range handler.metrics {
		_, err := fmt.Fprintf(w, "# HELP %s info msg\n# TYPE %s guage\n%s %f\n", k, k, k, v)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}

func (handler *ResourceConsumerHandler) bumpMetric(metric string, delta float64, duration time.Duration) {
	handler.metricsLock.Lock()
	if _, ok := handler.metrics[metric]; ok {
		handler.metrics[metric] += delta
	} else {
		handler.metrics[metric] = delta
	}
	handler.metricsLock.Unlock()

	time.Sleep(duration)

	handler.metricsLock.Lock()
	handler.metrics[metric] -= delta
	handler.metricsLock.Unlock()
}

func (handler *ResourceConsumerHandler) handleBumpMetric(w http.ResponseWriter, query url.Values) {
	metric := query.Get("metric")
	deltaString := query.Get("delta")
	durationSecString := query.Get("durationSec")

	if durationSecString == "" || metric == "" || deltaString == "" {
		http.Error(w, "durantionSec, metric, or delta missing", http.StatusBadRequest)
		return
	}

	durationSec, durationSecError := strconv.Atoi(durationSecString)
	delta, deltaError := strconv.ParseFloat(deltaString, 64)

	if durationSecError != nil || deltaError != nil {
		http.Error(w, "durantionSec or delta incorrect", http.StatusBadRequest)
		return
	}

	go handler.bumpMetric(metric, delta, time.Duration(durationSec)*time.Second)
	_, err := fmt.Fprintf(w, "/bump-metric\nmetric %s\ndelta %s\ndurationSec %d\n", metric, deltaString, durationSec)
	if err != nil {
		fmt.Println(err.Error())
	}
}
