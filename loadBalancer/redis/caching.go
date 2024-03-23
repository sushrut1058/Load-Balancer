package rCaching

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Cache struct {
	Status   int         `json:"StatusCode"`
	Header   http.Header `json:"Header"`
	Body     []byte      `json:"Body"`
	Validity time.Time   `json:"Validity"`
}

func GetCachedResponse(key string) (*Cache, error) {
	serialized_resp, err := redisClient.Get(ctx, key).Bytes()
	if err != nil {
		fmt.Println("Unable to fetch!\n", err)
		return nil, err
	}
	var CachedResp Cache
	err = json.Unmarshal(serialized_resp, &CachedResp)
	if err != nil {
		fmt.Println("Unable to deserialize!\n", err)
		return nil, err
	}
	return &CachedResp, nil

}

func SetCache(key string, body []byte, response *http.Response) {

	resp := &Cache{
		Status:   response.StatusCode,
		Header:   response.Header.Clone(),
		Body:     body,
		Validity: time.Now().Add(30 * time.Minute),
	}

	serialized_resp, err := json.Marshal(resp)
	if err != nil {
		fmt.Println("Parsing of Cache object in JSON has failed!\n", err)
		return
	}

	redisClient.Set(ctx, key, serialized_resp, 30*time.Minute)

}
