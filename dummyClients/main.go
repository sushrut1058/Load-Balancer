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
		}
		// io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
		time.Sleep(2000 * time.Millisecond)
	}
}
