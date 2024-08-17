package requests

import (
	"fmt"
	"net/http"
	"reverse-proxy/global"
	"sync"
)

var resourceMutex sync.Mutex // r/w on Capacity array, UrlIndex
var UrlIndex uint32
var counter int = 0

func getUrl() (string, uint32) {
	resourceMutex.Lock()
	defer resourceMutex.Unlock()
	localIndex := UrlIndex % uint32(global.NServers)
	index_ := DistributionStrategy(localIndex)
	return global.Servers[index_].URL, index_
}

func DistributionStrategy(index uint32) uint32 {
	return LoadBalancingStrategy(index)
}

func ReleaseResource(index uint32) {
	resourceMutex.Lock()
	defer resourceMutex.Unlock()
	global.CurrentCapacity[index] += 1
}

func copyHeaders(newHeader http.Header, rHeader http.Header) {
	for key, values := range rHeader {
		for _, value := range values {
			newHeader.Add(key, value)
		}
	}
}

func isPresentInJSON(a string, b []interface{}) bool {
	for _, url := range b {
		fmt.Println("debug", a, b)
		if a == url {
			fmt.Println("debug return TRUE")
			return true
		}
	}
	fmt.Println("debug return FALSE")
	return false
}
