package worker

import (
	"fmt"
	"reverse-proxy/global"
	"reverse-proxy/requests"
)

var workerChannel chan int

func StartWorkerPool() {
	requests.RequestChannel = make(chan requests.RequestHandle, 10000)
	workerChannel = make(chan int, global.MaxWorkerCount)
	for i := 0; i < global.MaxWorkerCount; i++ {
		workerChannel <- i
	}
	go TriggerWorkers()
}

func TriggerWorkers() {
	for requestHandle := range requests.RequestChannel {
		worker := <-workerChannel
		go Do(worker, requestHandle)
	}
}

func Do(i int, requestHandle requests.RequestHandle) {
	fmt.Println("Request picked by worker:", i)
	requestHandle.SendRequestAndForwardResponse()
	fmt.Println("Request processed by worker:", i)
	workerChannel <- i
}
