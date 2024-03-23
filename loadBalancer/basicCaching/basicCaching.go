package basicCaching

import (
	"net/http"
	"sync"
	"time"
)

type Cache struct {
	Status   int
	Header   http.Header
	Body     []byte
	Validity time.Time
}

var CacheMap = make(map[string]*Cache)
var cacheMutex = &sync.Mutex{}

func SetCache(key string, body []byte, response *http.Response) {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	CacheMap[key] = &Cache{
		Status:   response.StatusCode,
		Header:   response.Header.Clone(),
		Body:     body,
		Validity: time.Now().Add(30 * time.Minute),
	}

}

func GetCachedResponse(key string) (*Cache, bool) {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	resp, exists := CacheMap[key]
	if !exists || time.Now().After(resp.Validity) {
		return nil, false
	}
	return resp, true
}
