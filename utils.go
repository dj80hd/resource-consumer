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
	"log"
	"os/exec"
	"strconv"
)

//ConsumeCPU starts external process consuming millcores of CPU for durationSec
func ConsumeCPU(millicores int, durationSec int) {
	log.Printf("/consume-cpu millicores: %v, durationSec: %v", millicores, durationSec)
	// creating new consume cpu process
	arg1 := fmt.Sprintf("-millicores=%d", millicores)
	arg2 := fmt.Sprintf("-duration-sec=%d", durationSec)
	consumeCPU := exec.Command("/consume-cpu", arg1, arg2)
	err := consumeCPU.Run()
	if err != nil {
		log.Printf(err.Error())
	}
}

//ConsumeMem starts external process consuming millcores of CPU for durationSec
func ConsumeMem(megabytes int, durationSec int) {
	log.Printf("/consume-mem megabytes: %v, durationSec: %v", megabytes, durationSec)
	megabytesString := strconv.Itoa(megabytes) + "M"
	durationSecString := strconv.Itoa(durationSec)
	// creating new consume memory process
	consumeMem := exec.Command("stress", "-m", "1", "--vm-bytes", megabytesString, "--vm-hang", "0", "-t", durationSecString)
	err := consumeMem.Run()
	if err != nil {
		log.Printf(err.Error())
	}
}

//ConsumeDisk creates a file of the specified size
func ConsumeDisk(gigabytes int, filename string) {
	var err error
	log.Printf("/consume-disk gigabytes: %v file: %s", gigabytes, filename)
	arg1 := fmt.Sprintf("of=%s", filename)
	arg2 := fmt.Sprintf("count=%d", gigabytes)
	if gigabytes > 0 {
		consumeDisk := exec.Command("dd", "if=/dev/zero", "bs=1073741824", arg1, arg2)
		err = consumeDisk.Run()
	} else {
		freeDisk := exec.Command("rm", filename)
		err = freeDisk.Run()
	}
	if err != nil {
		log.Printf(err.Error())
	}
}

//GetCurrentStatus is not implemented
func GetCurrentStatus() {
	log.Printf("GetCurrentStatus")
}
