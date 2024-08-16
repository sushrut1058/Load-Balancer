package workers

import (
	"fmt"
	"reverse-proxy/global"
)

var CurrentCount int
var workerPool chan int

func InitializeWorkerPool(maxWorkers int) {
	workerPool = make(chan int, maxWorkers)
	for i := 0; i < maxWorkers; i++ {
		fmt.Println("[init]adding worker:", i)
		workerPool <- i
	}
}

func Do(q *global.Queue) { //, h global.RequestHandle
	fmt.Println("----------------[worker.Do() start]-------------------")
	workerID := <-workerPool
	fmt.Println("Worker acquired. ID:", workerID)
	fmt.Println("Popping queue of size:", q.Size())

	requestHandle, popped := q.Pop()
	if !popped || requestHandle == nil {
		fmt.Println("[worker.Do()] Empty queue or unable to pop!")
		workerPool <- workerID
		fmt.Println("[worker.Do()] Worker with workerID:", workerID, "has been freed!")
		return
	}

	reqHandle := requestHandle.(global.RequestHandle)
	global.SendRequestAndForwardResponse(reqHandle.Writer, reqHandle.Request, reqHandle.Processed)

	requestHandle = nil
	popped = false
	workerPool <- workerID
	fmt.Println("----------------[worker.Do() end]-------------------")
}
