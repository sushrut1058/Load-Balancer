package worker

import (
	"fmt"
	"reverse-proxy/global"
)

var workerChannel chan int

func InitializeWorkerPool() {
	global.RequestChannel = make(chan global.RequestHandle, 10000)
	workerChannel = make(chan int, global.MaxWorkerCount)
	for i := 0; i < global.MaxWorkerCount; i++ {
		workerChannel <- i
	}
}

func TriggerWorkers() {
	for requestHandle := range global.RequestChannel {
		worker := <-workerChannel
		go Do(worker, requestHandle)
	}
}

func Do(i int, requestHandle global.RequestHandle) {
	fmt.Println("Request picked by worker:", i)
	global.SendRequestAndForwardResponse(requestHandle.Writer, requestHandle.Request, requestHandle.Processed)
	fmt.Println("Request processed by worker:", i)
	workerChannel <- i
}
