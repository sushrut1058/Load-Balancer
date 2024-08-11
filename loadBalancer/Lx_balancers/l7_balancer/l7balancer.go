package l7

import (
	"bytes"
	"fmt"
	"io"
	store "loadBalancer/store"
	"net/http"
	"net/url"
	"unsafe"
)

func L7_balancer(w http.ResponseWriter, r *http.Request) {
	curUrl := store.GetUrl(store.Data["strategy"].(string))
	req_key := r.URL.String()
	isCacheable := (r.Method == http.MethodGet || r.Method == http.MethodHead) && !store.IsPresentInJSON(req_key, store.Data["cache-ignore"].([]interface{}))
	fmt.Println("isCacheable returned", isCacheable)
	parsedCurUrl, _ := url.Parse(curUrl.String())
	parsedCurUrl.Path = r.URL.Path
	parsedCurUrl.RawQuery = r.URL.RawQuery

	if isCacheable {
		fmt.Println("inside isCacheable")
		cResp, exists := store.GetCachedResponse(req_key)
		fmt.Println("returns", exists)
		if exists {
			fmt.Println("[l7loadbalancer] --------------Sending Cached Response----------------")
			for key, values := range cResp.Header {
				for _, value := range values {
					w.Header().Add(key, value)
				}
			}
			fmt.Println("Size of cached resposne: ", unsafe.Sizeof(cResp.Body))
			w.WriteHeader(cResp.Status)
			w.Write(cResp.Body)
			return
		}
	}

	//creation of a new request object with copied header info
	req, error := http.NewRequest(r.Method, parsedCurUrl.String(), r.Body)
	if error != nil {
		http.Error(w, "[l7loadbalancer] Error CREATING request", http.StatusInternalServerError)
		return
	}
	req.Header = r.Header
	//creation of a new client, sending the request, expecting response
	proxy := &http.Client{}
	resp, err := proxy.Do(req)
	if err != nil {
		http.Error(w, "[l7loadbalancer] Error FORWARDING request", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("[l7loadbalancer] Error while reading response: %v", err)
	}

	if isCacheable {
		fmt.Println("setting cache")
		store.SetCache(req_key, bodyBytes, resp)
	}

	resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}
	w.WriteHeader(resp.StatusCode)

	io.Copy(w, resp.Body)
}
