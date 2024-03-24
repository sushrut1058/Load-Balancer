package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"loadBalancer/basicCaching"
	rCaching "loadBalancer/redis"
	types "loadBalancer/types"
	"net"
	"net/http"
	"net/url"
	"os"
	"sync/atomic"
)

var servers []interface{}
var port string
var index uint64
var data map[string]interface{}

var SetCache func(string, []byte, *http.Response)
var GetCachedResponse func(string) (*types.Cache, bool)

func getUrl(strategy string) *url.URL {
	localIndex := atomic.AddUint64(&index, 1)
	curUrl := servers[localIndex%uint64(len(servers))].(string)
	retVal, _ := url.Parse(curUrl)
	return retVal
}

func l7balancer(w http.ResponseWriter, r *http.Request) {
	curUrl := getUrl(data["strategy"].(string))

	cResp, exists := GetCachedResponse(curUrl.String())
	if exists {
		fmt.Println("Using Cached Response")
		for key, values := range cResp.Header {
			for _, value := range values {
				w.Header().Add(key, value)
			}
		}
		w.WriteHeader(cResp.Status)
		w.Write(cResp.Body)
		return
	}

	//creation of a new request object with copied header info
	req, error := http.NewRequest(r.Method, curUrl.String(), r.Body)
	if error != nil {
		http.Error(w, "Error CREATING request", http.StatusInternalServerError)
		return
	}
	req.Header = r.Header
	//creation of a new client, sending the request, expecting response
	proxy := &http.Client{}
	resp, err := proxy.Do(req)
	if err != nil {
		http.Error(w, "Error FORWARDING request", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("[l7loadbalancer] Error while reading response: %v", err)
	}

	SetCache(curUrl.String(), bodyBytes, resp)

	resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}
	w.WriteHeader(resp.StatusCode)

	io.Copy(w, resp.Body)
}

func l4balancer(conn net.Conn) {
	defer conn.Close()
	curUrl := getUrl(data["strategy"].(string))
	serverConn, err := net.Dial(data["proto"].(string), curUrl.String())
	if err != nil {
		fmt.Printf("Failed to connect to backend %s: %v", curUrl.String(), err)
		conn.Close()
		return
	}
	defer serverConn.Close()

	go io.Copy(serverConn, conn)
	io.Copy(conn, serverConn)

}

func readConfig() {
	file, err := os.ReadFile("conf.json")
	if err != nil {
		fmt.Printf("Error while reading json: %v\n", err)
	}
	err = json.Unmarshal(file, &data)
	if err != nil {
		fmt.Printf("Error while writing json to memory: %v\n", err)
	}

	//caching
	if data["caching"] == "baseline" {
		SetCache = basicCaching.SetCache
		GetCachedResponse = basicCaching.GetCachedResponse
	} else if data["caching"] == "redis" {
		SetCache = rCaching.SetCache
		GetCachedResponse = rCaching.GetCachedResponse
	} else {
		SetCache = func(string, []byte, *http.Response) {}
		GetCachedResponse = func(string) (*types.Cache, bool) { return nil, false }
	}

	//servers and port
	servers, _ = data["servers"].([]interface{})
	port = fmt.Sprintf(":%v", data["port"].(string))

}

func main() {
	readConfig()

	if data["level"] == "L7" {
		http.HandleFunc("/", l7balancer)
		http.ListenAndServe(port, nil)
	} else {
		listener, err := net.Listen(data["proto"].(string), port)
		if err != nil {
			fmt.Printf("Error creating listener: %v\n", err)
		}
		defer listener.Close()
		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Printf("Error accepting connection: %v\n", err)
				continue
			}
			go l4balancer(conn)
		}
	}
}
