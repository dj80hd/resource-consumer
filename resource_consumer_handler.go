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
	if req.URL.Path == "/Metrics" {
		handler.handleMetrics(w)
		return
	}
	if req.Method != "POST" {
		http.Error(w, "Bad", http.StatusBadRequest)
		return
	}
	// parsing POST request.data and URL data
	if err := req.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if req.URL.Path == "/ConsumeCPU" {
		handler.handleConsumeCPU(w, req.Form)
		return
	}
	if req.URL.Path == "/ConsumeMem" {
		handler.handleConsumeMem(w, req.Form)
		return
	}
	if req.URL.Path == "/ConsumeDisk" {
		handler.handleConsumeDisk(w, req.Form)
		return
	}
	// handle getCurrentStatus
	if req.URL.Path == "/GetCurrentStatus" {
		handler.handleGetCurrentStatus(w)
		return
	}
	// handle bumpMetric
	if req.URL.Path == "/BumpMetric" {
		handler.handleBumpMetric(w, req.Form)
		return
	}
	http.Error(w, fmt.Sprintf("%s: %s", "unknown", req.URL.Path), http.StatusNotFound)
}

func (handler *ResourceConsumerHandler) handleConsumeCPU(w http.ResponseWriter, query url.Values) {
	// getting string data for consumeCPU
	durationSecString := query.Get("durationSec")
	millicoresString := query.Get("millicores")
	if durationSecString == "" || millicoresString == "" {
		http.Error(w, "not", http.StatusBadRequest)
		return
	}

	// convert data (strings to ints) for consumeCPU
	durationSec, durationSecError := strconv.Atoi(durationSecString)
	millicores, millicoresError := strconv.Atoi(millicoresString)
	if durationSecError != nil || millicoresError != nil {
		http.Error(w, "incorrect", http.StatusBadRequest)
		return
	}

	go ConsumeCPU(millicores, durationSec)
	fmt.Fprintln(w, "ConsumeCPU")
	fmt.Fprintln(w, "millicores ", millicores)
	fmt.Fprintln(w, "durationSec ", durationSec)
}

func (handler *ResourceConsumerHandler) handleConsumeDisk(w http.ResponseWriter, query url.Values) {
	filename := query.Get("filename")
	gigabytesString := query.Get("gigabytes")
	if filename == "" || gigabytesString == "" {
		http.Error(w, "filename or gigabytes missing", http.StatusBadRequest)
		return
	}
	fmt.Fprintln(w, "ConsumeDisk")
	fmt.Fprintln(w, "gigabytes ", gigabytesString)
	fmt.Fprintln(w, "filename ", filename)
}

func (handler *ResourceConsumerHandler) handleConsumeMem(w http.ResponseWriter, query url.Values) {
	// getting string data for consumeMem
	durationSecString := query.Get("durationSec")
	megabytesString := query.Get("megabytes")
	if durationSecString == "" || megabytesString == "" {
		http.Error(w, "durantionSecString or megabytesString missing", http.StatusBadRequest)
		return
	}

	// convert data (strings to ints) for consumeMem
	durationSec, durationSecError := strconv.Atoi(durationSecString)
	megabytes, megabytesError := strconv.Atoi(megabytesString)
	if durationSecError != nil || megabytesError != nil {
		http.Error(w, "durantionSecError or megabytesError incorrect", http.StatusBadRequest)
		return
	}

	go ConsumeMem(megabytes, durationSec)
	fmt.Fprintln(w, "ConsumeMem")
	fmt.Fprintln(w, "megabytes ", megabytes)
	fmt.Fprintln(w, "durationSec ", durationSec)
}

func (handler *ResourceConsumerHandler) handleGetCurrentStatus(w http.ResponseWriter) {
	GetCurrentStatus()
	fmt.Fprintln(w, "Warning: not implemented!")
	fmt.Fprint(w, "GetCurrentStatus")
}

func (handler *ResourceConsumerHandler) handleMetrics(w http.ResponseWriter) {
	handler.metricsLock.Lock()
	defer handler.metricsLock.Unlock()
	for k, v := range handler.metrics {
		fmt.Fprintf(w, "# HELP %s info message.\n", k)
		fmt.Fprintf(w, "# TYPE %s gauge\n", k)
		fmt.Fprintf(w, "%s %f\n", k, v)
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
	// getting string data for handleBumpMetric
	metric := query.Get("metric")
	deltaString := query.Get("delta")
	durationSecString := query.Get("durationSec")
	if durationSecString == "" || metric == "" || deltaString == "" {
		http.Error(w, "durantionSecString, metric, or deltaString missing", http.StatusBadRequest)

		return
	}

	// convert data (strings to ints/floats) for handleBumpMetric
	durationSec, durationSecError := strconv.Atoi(durationSecString)
	delta, deltaError := strconv.ParseFloat(deltaString, 64)
	if durationSecError != nil || deltaError != nil {
		http.Error(w, "durantionSec or deltaString incorrect", http.StatusBadRequest)
		return
	}

	go handler.bumpMetric(metric, delta, time.Duration(durationSec)*time.Second)
	fmt.Fprintln(w, "BumpMetric")

	fmt.Fprintln(w, "metric ", metric)
	fmt.Fprintln(w, "delta ", delta)
	fmt.Fprintln(w, "durationSec ", durationSec)
}
