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
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"
)

var port = flag.Int("port", 8080, "Port number.")
var secs = flag.Int("secs", -1, "seconds to remain running; -1 is forever")

func main() {
	flag.Parse()

	if *secs > 1 {
		go func() {
			time.Sleep(time.Duration(*secs) * time.Second)
			log.Fatalf("exit after %d seconds", *secs)
		}()
	}

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), NewResourceConsumerHandler()))
}
