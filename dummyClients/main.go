package main

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

func main() {
	fmt.Println("Started!")
	for i := 0; i < 60; i++ {
		resp, err := http.Get("http://localhost:8080")
		start:=time.Now()
		if err != nil {
			fmt.Printf("Error in request: %v\n", err)
			continue
		}
		fmt.Printf("Request sent. Status Code: %d\n", resp.StatusCode)
		body, error := io.ReadAll(resp.Body)
		if error != nil {
			fmt.Printf("Error while reading: %v\n", error)
		} else {
			resp_string := string(body)
			fmt.Printf("[RESPONSE]: %v\n", resp_string)
			fmt.Printf("----------[RESPONSE TIME]-------------: %v\n", time.Now().Sub(start))
		}
		// io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
		time.Sleep(2000 * time.Millisecond)
	}
}
