package rCaching

import (
	"fmt"
	"net/http"
	"time"
)

type Cache struct {
	StatusCode int         `json:"StatusCode"`
	Header     http.Header `json:"Header"`
	Body       []byte      `json:"Body"`
	Validity   time.Time   `json:"Validity"`
}

func GetCachedResponse(key string) (Cache, bool) {
	key_header := fmt.Sprintf("%v_Header", key)
	key_statusCode := fmt.Sprintf("%v_StatusCode", key)

}

func SetCache(key string, body []byte, response *http.Response) {

}
