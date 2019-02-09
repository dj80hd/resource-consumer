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
	"time"
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

//ConsumeMem starts external process consuming megabytes of MEM for durationSec
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

//ConsumeDisk creates file gigabtyes in size, deletes it after durationSec
func ConsumeDisk(gigabytes int, durationSec int, filename string) {
	var err error
	log.Printf("/consume-disk gigabytes: %v filename: %s durationSec: %v",
		gigabytes, filename, durationSec)

	if gigabytes <= 0 {
		return
	}

	consumeDisk := exec.Command(
		"dd", "if=/dev/zero", "bs=1048576",
		fmt.Sprintf("of=%s", filename),
		fmt.Sprintf("count=%d", gigabytes*1024))
	freeDisk := exec.Command("rm", filename)

	err = consumeDisk.Run()
	if err != nil {
		log.Printf(err.Error())
	}

	go func(durationSec int, cmd *exec.Cmd) {
		time.Sleep(time.Duration(durationSec * 1e9))
		err := cmd.Run()
		if err != nil {
			log.Printf(err.Error())
		}
	}(durationSec, freeDisk)
}
