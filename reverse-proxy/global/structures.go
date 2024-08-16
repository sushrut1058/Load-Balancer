package global

import (
	"net/http"
	"sync"
)

type Resource struct {
	URL      string
	Capacity float64
}

type Queue struct {
	items []interface{}
	mtx   sync.Mutex
}

type RequestHandle struct {
	Request   *http.Request
	Writer    http.ResponseWriter
	Processed *chan bool
}
