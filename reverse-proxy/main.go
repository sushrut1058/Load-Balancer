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
	newRequestHandle := global.RequestHandle{Request: r, Writer: w}
	global.RequestQueue.Push(newRequestHandle)
	fmt.Println("calling worker")
	workers.Do(&global.RequestQueue) //, newRequestHandle)
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

	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
