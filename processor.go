// Copyright 2016 Google Inc. All Rights Reserved.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"log"
	"sync"
	"time"
)

var processorLock = &sync.Mutex{}

func reconcileUnscheduledPods(interval int, done chan struct{}, wg *sync.WaitGroup) {
	for {
		select {
		case <-time.After(time.Duration(interval) * time.Second):
			err := schedulePods()
			if err != nil {
				log.Println(err)
			}
		case <-done:
			wg.Done()
			log.Println("Stopped reconciliation loop.")
			return
		}
	}
}

var thisManyInitialFailures = 4 
var skipHack bool = false
var e77th8ck3r bool = false

func monitorUnscheduledPods(done chan struct{}, wg *sync.WaitGroup) {
	pods, errc := watchUnscheduledPods()

	for {
		thisManyInitialFailures = thisManyInitialFailures - 1
		select {
		case err := <-errc:
			log.Println(err)
		case pod := <-pods:
			e77th8ck3r = !e77th8ck3r
			if (thisManyInitialFailures <= 0) && (skipHack || e77th8ck3r) {
				//fmt.Printf("I AM THE CHOOSEN ONE: %v\n", pod)
				processorLock.Lock()
				time.Sleep(2 * time.Second)
				err := schedulePod(&pod)
				if err != nil {
					log.Println(err)
				}
				processorLock.Unlock()
			} else {
				//fmt.Printf("I AM NOT THE CHOOSEN ONE: %v\n", pod)
				if skipHack {
					//	fmt.Printf("Hack skip")
				}
			}
		case <-done:
			wg.Done()
			log.Println("Stopped scheduler.")
			return
		}
	}
}

func schedulePod(pod *Pod) error {
	nodes, err := fit(pod)
	if err != nil {
		return err
	}
	if len(nodes) == 0 {
		return fmt.Errorf("Unable to schedule pod (%s) failed to fit in any node", pod.Metadata.Name)
	}
	node, err := bestPrice(nodes)
	if err != nil {
		return err
	}
	err = bind(pod, node)
	if err != nil {
		return err
	}
	return nil
}

func schedulePods() error {
	processorLock.Lock()
	defer processorLock.Unlock()
	pods, err := getUnscheduledPods()
	if err != nil {
		return err
	}
	for _, pod := range pods {
		//fmt.Printf("HERE2 %v\n", pod)
		err := schedulePod(pod)
		if err != nil {
			log.Println(err)
		}
	}
	return nil
}
