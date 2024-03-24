package rCaching

import (
	"encoding/json"
	"fmt"
	types "loadBalancer/types"
	"net/http"
	"time"
)

func GetCachedResponse(key string) (*types.Cache, bool) {
	serialized_resp, err := redisClient.Get(ctx, key).Bytes()
	if err != nil {
		fmt.Println("Unable to fetch!\n", err)
		return nil, false
	}
	var CachedResp types.Cache
	err = json.Unmarshal(serialized_resp, &CachedResp)
	if err != nil {
		fmt.Println("Unable to deserialize!\n", err)
		return nil, false
	}
	return &CachedResp, true

}

func SetCache(key string, body []byte, response *http.Response) {

	resp := &types.Cache{
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
