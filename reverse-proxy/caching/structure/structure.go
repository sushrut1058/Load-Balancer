package structure

import (
	"net/http"
	"time"
)

type Cache struct {
	Status   int         `json:"StatusCode"`
	Header   http.Header `json:"Header"`
	Body     []byte      `json:"Body"`
	Validity time.Time   `json:"Validity"`
}
