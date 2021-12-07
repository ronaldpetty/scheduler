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
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

const schedulerName = "default-scheduler"

func checkError(err error) {
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	go func() {
	//	log.Println("Listening healthz over TLS")
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "")
		})
		err := http.ListenAndServeTLS("localhost:10259", "server.crt", "server.key", nil)
		checkError(err)
	}()

	//log.Println("Starting custom scheduler...")

	doneChan := make(chan struct{})
	var wg sync.WaitGroup

	wg.Add(1)
	go monitorUnscheduledPods(doneChan, &wg)

	wg.Add(1)
	go reconcileUnscheduledPods(30, doneChan, &wg)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	for {
		select {
		case <-signalChan:
			log.Printf("Shutdown signal received, exiting...")
			close(doneChan)
			wg.Wait()
			os.Exit(0)
		}
	}

}
