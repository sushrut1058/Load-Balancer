package inMemory

import (
	"net/http"
	"reverse-proxy/caching/structure"
	"sync"
	"time"
)

var CacheMap = make(map[string]*structure.Cache)
var cacheMutex = &sync.Mutex{}

func SetCache(key string, body []byte, response *http.Response) {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	CacheMap[key] = &structure.Cache{
		Status:   response.StatusCode,
		Header:   response.Header.Clone(),
		Body:     body,
		Validity: time.Now().Add(30 * time.Minute),
	}

}

func GetCachedResponse(key string) (*structure.Cache, bool) {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	resp, exists := CacheMap[key]
	if !exists || time.Now().After(resp.Validity) {
		delete(CacheMap, key)
		return nil, false
	}
	return resp, true
}
