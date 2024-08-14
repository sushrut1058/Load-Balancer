package global

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
)

var resourceMutex sync.Mutex // r/w on Capacity array, UrlIndex

func getUrl() (string, uint32) {
	resourceMutex.Lock()
	defer resourceMutex.Unlock()
	index := UrlIndex % uint32(NServers)
	DistributionStrategy(index)
	return Servers[index].URL, index
}

func DistributionStrategy(index uint32) {
	if CurrentCapacity[index] > 0 {
		CurrentCapacity[index] -= 1
	} else {
		CurrentCapacity[index] = int(Servers[index].Capacity)
		UrlIndex += 1
	}
}

func ReleaseResource(index uint32) {
	resourceMutex.Lock()
	defer resourceMutex.Unlock()
	CurrentCapacity[index] += 1
}

func sendRequest(r *http.Request) *http.Response {
	fmt.Println("[global.sendRequest()] [Inside] getUrl")
	url_string, index := getUrl()
	url_string = url_string + r.RequestURI
	url, _ := url.Parse(url_string)
	fmt.Println("[global.sendRequest()] ExternalSendRequest calling . . .")
	resp := ExternalSendRequest(url, r)
	if resp == nil {
		fmt.Println("Response is nil!!!!")
		return nil
	}
	body, _ := io.ReadAll(resp.Body)
	fmt.Println(string(body))
	fmt.Println("[global.sendRequest()] Releasing resources . . .")
	ReleaseResource(index)
	if resp == nil {
		return nil
	}
	resp.Body = io.NopCloser(bytes.NewBuffer(body))
	return resp
}

/* START - Might shift to layer specific code*/
func ExternalSendRequest(url *url.URL, r *http.Request) *http.Response {
	newRequest, err := http.NewRequest(r.Method, url.String(), r.Body)
	if err != nil {
		fmt.Println("Couldn't create a new request object!")
		return nil
	}
	copyHeaders(newRequest.Header, r.Header)
	client := &http.Client{}
	resp, err := client.Do(newRequest)
	if err != nil {
		fmt.Println("Couldn't connect to server. Error:", err)
		return nil
	}
	return resp
}

func copyHeaders(newHeader http.Header, rHeader http.Header) {
	for key, values := range rHeader {
		for _, value := range values {
			newHeader.Add(key, value)
		}
	}
}

/* END - Might shift to layer specific code*/

func forwardResponse(resp *http.Response, w http.ResponseWriter) {
	copyHeaders(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func SendRequestAndForwardResponse(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[global.SendRequestAndForwardResponse()] [Inside]")
	resp := sendRequest(r)
	defer resp.Body.Close()
	forwardResponse(resp, w)
}
