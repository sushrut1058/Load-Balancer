package global

import (
	"fmt"
	"sort"
)

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
		ServerIndexMap[servers[i].URL] = i
	}
	NServers = len(servers)
	Servers = servers
}
