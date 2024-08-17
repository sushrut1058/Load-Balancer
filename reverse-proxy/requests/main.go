package requests

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"reverse-proxy/caching"
	"reverse-proxy/global"
)

var RequestChannel chan RequestHandle

type HTTPRequestHandle struct {
	Request   *http.Request
	Writer    http.ResponseWriter
	Processed *chan bool
}

type TCPRequestHandle struct {
	Conn net.Conn
}

type RequestHandle interface {
	SendRequestAndForwardResponse()
}

func (handle *HTTPRequestHandle) SendRequestAndForwardResponse() {

	r, w, processed := handle.Request, handle.Writer, handle.Processed
	//backend server selection
	url_string, index := getUrl()
	url_string = url_string + r.RequestURI
	url, err := url.Parse(url_string)
	if err != nil {
		fmt.Println("Error occurred while parsing url. Error:", err)
		return
	}

	isCacheable := (r.Method == http.MethodGet || r.Method == http.MethodHead) && !isPresentInJSON(url_string, global.Data["cache-ignore"].([]interface{}))
	if isCacheable {
		fmt.Println("Response is cacheable")
		cResp, exists := caching.GetCachedResponse(r.URL.String())
		if exists {
			fmt.Println("Response already exists in cache")
			copyHeaders(w.Header(), cResp.Header)
			w.WriteHeader(cResp.Status)
			w.Write(cResp.Body)
			*processed <- true
			return
		} else {
			fmt.Println("Response was not found")
		}
	}
	//sending request
	newRequest, err := http.NewRequest(r.Method, url.String(), r.Body)
	if err != nil {
		fmt.Println("Couldn't create a new request object!")
		return
	}
	copyHeaders(newRequest.Header, r.Header)

	client := &http.Client{}
	resp, err := client.Do(newRequest)
	if err != nil {
		fmt.Println("Couldn't connect to server. Error:", err)
		return
	}
	defer resp.Body.Close()
	defer ReleaseResource(index)

	fmt.Println("StatusCode:", resp.StatusCode)
	if isCacheable {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Trouble converting response body to bytes. Error:", err)
		}
		fmt.Println("Setting Cache for body:", string(bodyBytes))
		caching.SetCache(r.URL.String(), bodyBytes, resp)
		resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	}

	copyHeaders(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)

	*processed <- true
}

func (handle *TCPRequestHandle) SendRequestAndForwardResponse() {
	conn := handle.Conn
	defer conn.Close()

	url_string, _ := getUrl()

	serverConn, err := net.Dial(global.Data["proto"].(string), url_string)
	if err != nil {
		fmt.Printf("[l4balancer] Failed to connect to backend %s: %v", url_string, err)
		conn.Close()
		return
	}
	defer serverConn.Close()

	go io.Copy(serverConn, conn)
	io.Copy(conn, serverConn)
}
