package workers

import (
	"fmt"
	"reverse-proxy/global"
	"sync"
)

var CurrentCount int
var workerCountMutex sync.Mutex

func Do(q *global.Queue) {
	for q.Size() != 0 {
		var workerid int
		{
			workerCountMutex.Lock()
			workerid = CurrentCount
			fmt.Println("[worker.Do()] Upper Mutex LOCKED", workerid)
			fmt.Println("[worker.Do()] CurrentCount:", CurrentCount)
			if CurrentCount >= global.MaxWorkerCount {
				workerCountMutex.Unlock()
				break
			} else {
				CurrentCount++
			}
			workerCountMutex.Unlock()
			fmt.Println("[worker.Do()] Upper Mutex UNLOCKED", workerid)
		}
		requestHandle, popped := q.Pop()
		if !popped || requestHandle == nil {
			fmt.Println("[worker.Do()] Empty queue or unable to pop!")
			continue
		}
		reqHandle := requestHandle.(global.RequestHandle)
		global.SendRequestAndForwardResponse(reqHandle.Writer, reqHandle.Request)
		{
			workerCountMutex.Lock()
			fmt.Println("[worker.Do()] Bottom Mutex LOCKED", workerid)
			CurrentCount--
			workerCountMutex.Unlock()
			fmt.Println("[worker.Do()] Bottom Mutex UNLOCKED", workerid)
		}
		fmt.Println("One cycle done, for workerid:", workerid)
		fmt.Println("Current status of queue:", q)
		fmt.Println("----------------!!!-------------------")
	}
}
