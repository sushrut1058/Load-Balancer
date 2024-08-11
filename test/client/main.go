package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func main() {
	fmt.Println("Started!")
	file, err := os.ReadFile("body.txt")
	if err != nil {
		fmt.Println("Error reading file")
	}

	// fmt.Println(file)

	client := &http.Client{}
	request, err := http.NewRequest("GET", "http://localhost:8080/api", nil)
	if err != nil {
		fmt.Println("", request, err)
	}
	avg := 0
	k := 10
	for i := 0; i < k; i++ {
		resp, err := client.Do(request)
		start := time.Now()
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
			va := time.Since(start)
			fmt.Printf("----------[GET RESPONSE TIME]-------------: %v\n", va)
			avg += int(va)
		}

		// io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
		time.Sleep(50 * time.Millisecond)
	}
	fmt.Println("Average response time for GET requests: ", avg/k)
	fmt.Println("--------------------------------------------------------")
	request, err = http.NewRequest("POST", "http://localhost:8080", bytes.NewReader(file))
	if err != nil {
		fmt.Println("", request, err)
	}
	avg = 0
	k = 10
	for i := 0; i < k; i++ {
		resp, err := client.Do(request)
		start := time.Now()
		if err != nil {
			fmt.Printf("Error in request: %v\n", err)
			continue
		}
		fmt.Printf("Request sent. Status Code: %d\n", resp.StatusCode)
		_, error := io.ReadAll(resp.Body)
		if error != nil {
			fmt.Printf("Error while reading: %v\n", error)
		} else {
			// resp_string := string(body)
			// fmt.Printf("[RESPONSE]: %v\n", resp_string)
			va := time.Since(start)
			fmt.Printf("----------[POST RESPONSE TIME]-------------: %v\n", va)
			avg += int(va)
		}

		// io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
		time.Sleep(50 * time.Millisecond)
	}
}
