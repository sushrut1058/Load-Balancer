package workers

import (
	"fmt"
	"reverse-proxy/global"
	"sync"
	"time"
)

var CurrentCount int
var workerCountMutex sync.Mutex

func Do(q *global.Queue) { //, h global.RequestHandle
	fmt.Println("----------------[worker.Do() start]-------------------")
	for q.Size() != 0 {
		// var workerid int
		{
			workerCountMutex.Lock()
			// workerid = CurrentCount
			if CurrentCount >= global.MaxWorkerCount {
				workerCountMutex.Unlock()
				time.Sleep(100 * time.Millisecond)
				continue
			} else {
				fmt.Println("Acquiring free worker spot")
				// if h != (global.RequestHandle{}) {
				// 	q.Push(h)
				// }
				CurrentCount++
			}
			workerCountMutex.Unlock()
		}
		fmt.Println("Popping queue of size:", q.Size())
		requestHandle, popped := q.Pop()
		if !popped || requestHandle == nil {
			fmt.Println("[worker.Do()] Empty queue or unable to pop!")
			continue
		}
		reqHandle := requestHandle.(global.RequestHandle)
		global.SendRequestAndForwardResponse(reqHandle.Writer, reqHandle.Request)
		{
			workerCountMutex.Lock()
			CurrentCount--
			workerCountMutex.Unlock()
		}
		requestHandle = nil
		popped = false
		fmt.Println("----------------[worker.Do() end]-------------------")
	}
}
