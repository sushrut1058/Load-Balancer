package store

import (
	"fmt"
	"loadBalancer/types"
	"net/http"
	"net/url"
	"sync/atomic"
)

var Servers []interface{}
var Port string
var Index uint64
var Data map[string]interface{}

var SetCache func(string, []byte, *http.Response)
var GetCachedResponse func(string) (*types.Cache, bool)

// Function to get URL based on the current strategy
func GetUrl(strategy string) *url.URL {
	localIndex := atomic.AddUint64(&Index, 1)
	curUrl := Servers[localIndex%uint64(len(Servers))].(string)
	retVal, _ := url.Parse(curUrl)
	return retVal
}

func IsPresentInJSON(a string, b []interface{}) bool {
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
