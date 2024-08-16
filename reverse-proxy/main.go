package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	global "reverse-proxy/global"
	"reverse-proxy/workers"
)

func handler(w http.ResponseWriter, r *http.Request) {
	done := make(chan bool)
	newRequestHandle := global.RequestHandle{Request: r, Writer: w, Processed: &done}
	global.RequestQueue.Push(newRequestHandle)
	fmt.Println("calling worker")
	go workers.Do(&global.RequestQueue) //, newRequestHandle)
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

	readConfiguration()

	fmt.Println(global.MaxWorkerCount)
	workers.InitializeWorkerPool(global.MaxWorkerCount)

	http.HandleFunc("/", handler)
	fmt.Println("Handler added")
	fmt.Println("Listening...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error while listening. error:", err)
	}

}
