package global

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
)

func Preprocessing() {
	readConfiguration()
}

func InitServerMap(serversJson map[string]interface{}) {
	var servers []Resource
	ServerIndexMap = make(map[string]int)
	for key, value := range serversJson {
		fmt.Println("starting")
		obj := Resource{URL: key, Capacity: value.(float64)}
		servers = append(servers, obj)
	}
	sort.SliceStable(servers, func(i, j int) bool {
		return servers[i].URL > servers[j].URL
	})
	for i := range servers {
		CurrentCapacity = append(CurrentCapacity, int(servers[i].Capacity))
		TotalCapacity = append(TotalCapacity, int(servers[i].Capacity))
		ServerIndexMap[servers[i].URL] = i
	}
	NServers = len(servers)
	Servers = servers
}

func readConfiguration() {
	file, err := os.ReadFile("config.json")
	if err != nil {
		fmt.Println("[main readConfiguration] Error reading file. Error:", err)
		return
	}
	if err = json.Unmarshal(file, &Data); err != nil {
		fmt.Println("[main readConfiguration] Error unmarshaling file into struct")
		return
	}
	MaxWorkerCount = int(Data["maxWorkers"].(float64))
	InitServerMap(Data["servers"].(map[string]interface{}))
	fmt.Println("CurrentCapacity:", CurrentCapacity)
	fmt.Println("Servers:", Servers)
}
