package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync/atomic"
)

var servers = []string{"http://localhost:8081", "http://localhost:8082", "http://localhost:8083"}
var index uint64

func getUrl() *url.URL {
	localIndex := atomic.AddUint64(&index, 1)
	curUrl := servers[localIndex%uint64(len(servers))]
	retVal, _ := url.Parse(curUrl)
	return retVal

}

func loadBalancer(w http.ResponseWriter, r *http.Request) {
	curUrl := getUrl()
	proxy := httputil.NewSingleHostReverseProxy(curUrl)
	proxy.ErrorHandler = func(writer http.ResponseWriter, request *http.Request, e error) {
		fmt.Printf("Error sending request to %v, due to %v\n", curUrl, e)
		http.Error(writer, "Service unavailable\n", http.StatusServiceUnavailable)
	}
	fmt.Printf("Forwarding request to %v!\n", curUrl)
	proxy.ServeHTTP(w, r)
}

func main() {
	http.HandleFunc("/", loadBalancer)
	http.ListenAndServe(":8080", nil)
	fmt.Println("Listening...")
}
