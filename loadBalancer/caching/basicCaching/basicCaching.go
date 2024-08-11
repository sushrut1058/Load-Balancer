package basicCaching

import (
	types "loadBalancer/types"
	"net/http"
	"sync"
	"time"
)

var CacheMap = make(map[string]*types.Cache)
var cacheMutex = &sync.Mutex{}

func SetCache(key string, body []byte, response *http.Response) {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	CacheMap[key] = &types.Cache{
		Status:   response.StatusCode,
		Header:   response.Header.Clone(),
		Body:     body,
		Validity: time.Now().Add(30 * time.Minute),
	}

}

func GetCachedResponse(key string) (*types.Cache, bool) {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	resp, exists := CacheMap[key]
	if !exists || time.Now().After(resp.Validity) {
		delete(CacheMap, key)
		return nil, false
	}
	return resp, true
}
