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
	fmt.Println("Fing around with vars")
	defer resourceMutex.Unlock()
	localIndex := UrlIndex % uint32(NServers)
	index_ := DistributionStrategy(localIndex)
	fmt.Println(UrlIndex)
	fmt.Println("Stopping Fing around with vars")
	return Servers[index_].URL, index_
}

func DistributionStrategy(index uint32) uint32 {
	return weightedRoundRobin(index)
}

func weightedRoundRobin(index uint32) uint32 {
	if TotalCapacity[index] == counter {
		counter = 0
		UrlIndex++
	}
	index = UrlIndex % uint32(NServers)
	counter++
	return index
}

func ReleaseResource(index uint32) {
	resourceMutex.Lock()
	defer resourceMutex.Unlock()
	CurrentCapacity[index] += 1
}

func sendRequest(r *http.Request) ([]byte, *http.Response, uint32) {
	fmt.Println("[global.sendRequest()] [Inside] getUrl")
	url_string, index := getUrl()
	url_string = url_string + r.RequestURI
	url, _ := url.Parse(url_string)
	fmt.Println("[global.sendRequest()] ExternalSendRequest calling . . .")
	bodyBytes, resp := ExternalSendRequest(url, r)
	if bodyBytes == nil {
		fmt.Println("Response bytes are nil!!!!")
		return nil, &http.Response{}, 0
	}
	fmt.Println("[sendRequest] Body:", string(bodyBytes))
	fmt.Println("[global.sendRequest()] Releasing resources . . .")
	return bodyBytes, resp, index
}

/* START - Might shift to layer specific code*/
func ExternalSendRequest(url *url.URL, r *http.Request) ([]byte, *http.Response) {
	newRequest, err := http.NewRequest(r.Method, url.String(), r.Body)
	if err != nil {
		fmt.Println("Couldn't create a new request object!")
		return nil, &http.Response{}
	}
	// copyHeaders(newRequest.Header, r.Header)
	for key, values := range r.Header {
		for _, value := range values {
			newRequest.Header.Add(key, value)
		}
	}
	client := &http.Client{}
	resp, err := client.Do(newRequest)

	if err != nil {
		fmt.Println("Couldn't connect to server. Error:", err)
		return nil, &http.Response{}
	}
	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body in bytes. Error:", err)
	}
	return bodyBytes, resp
}

// func copyHeaders(newHeader http.Header, rHeader http.Header) {
// 	for key, values := range rHeader {
// 		for _, value := range values {
// 			newHeader.Add(key, value)
// 		}
// 	}
// }

/* END - Might shift to layer specific code*/

func forwardResponse(bodyBytes []byte, resp *http.Response, w http.ResponseWriter) {
	fmt.Println("StatusCode:", resp.StatusCode)
	// copyHeaders(w.Header(), respHeader)
	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}
	w.WriteHeader(resp.StatusCode)
	fmt.Println("Header:", w.Header())

	fmt.Println("[FORWARD-RESPONSE] BODY:", string(bodyBytes))

	io.Copy(w, io.NopCloser(bytes.NewBuffer(bodyBytes)))
}

func SendRequestAndForwardResponse(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[global.SendRequestAndForwardResponse()] [Inside]")
	bodyBytes, resp, index := sendRequest(r)
	defer ReleaseResource(index)
	fmt.Println("[sendandforwardresponse] after sendrequest, printing body:", string(bodyBytes))
	fmt.Println("[sendandforwardresponse] after sendrequest, printing respHeader:", resp.Header)

	forwardResponse(bodyBytes, resp, w)

}
