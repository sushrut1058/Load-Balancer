package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	global "reverse-proxy/global"
	"reverse-proxy/worker"
	"time"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("request reached")
	done := make(chan bool)
	newRequestHandle := global.RequestHandle{Request: r, Writer: w, Processed: &done}
	// global.RequestQueue.Push(newRequestHandle)
	global.RequestChannel <- newRequestHandle
	<-done
}

func readConfiguration() {
	file, err := os.ReadFile("config.json")
	if err != nil {
		fmt.Println("[main readConfiguration] Error reading file. Error:", err)
		return
	}
	if err = json.Unmarshal(file, &global.Data); err != nil {
		fmt.Println("[main readConfiguration] Error unmarshaling file into struct")
		return
	}
	// fmt.Println(global.Data)
	global.MaxWorkerCount = int(global.Data["maxWorkers"].(float64))
	global.InitServerMap(global.Data["servers"].(map[string]interface{}))
	fmt.Println("CurrentCapacity:", global.CurrentCapacity)
	fmt.Println("Servers:", global.Servers)
}

func main() {
	fmt.Println("Starting . . .")
	// go func() {
	// 	for {
	// 		fmt.Printf("Number of goroutines: %d\n", runtime.NumGoroutine())
	// 		time.Sleep(10 * time.Millisecond)
	// 	}
	// }()

	readConfiguration()

	worker.InitializeWorkerPool()
	go worker.TriggerWorkers()

	http.HandleFunc("/", handler)
	fmt.Println("Handler added")
	fmt.Println("Listening...")
	time.Sleep(2 * time.Second)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error while listening. error:", err)
	}

}
