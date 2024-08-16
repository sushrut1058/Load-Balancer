package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

func main() {
	fmt.Println("Started!")

	var wg sync.WaitGroup

	client := &http.Client{
		Transport: &http.Transport{
			DisableKeepAlives: true, // Crucial for testing, each request to be considered a different client.
		},
	}
	// go func() {
	// 	for {
	// 		fmt.Printf("Number of goroutines: %d\n", runtime.NumGoroutine())
	// 		time.Sleep(10 * time.Millisecond)
	// 	}
	// }()

	k := 10 // Number of concurrent requests
	start := time.Now()
	for i := 0; i < k; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			// Create the request body containing the ID
			body := bytes.NewBufferString(fmt.Sprintf("Request ID: %d", id))
			request, err := http.NewRequest("POST", "http://localhost:8080/api", body)
			if err != nil {
				fmt.Printf("Error creating request: %v\n", err)
				return
			}

			// Set the content type to indicate the body is plain text
			request.Header.Set("Content-Type", "text/plain")

			resp, err := client.Do(request)

			if err != nil {
				fmt.Printf("Error in request ID %d: %v\n", id, err)
				return
			}

			// Read and print the response
			respBody, err := io.ReadAll(resp.Body)
			if err != nil {
				fmt.Printf("Error while reading response ID %d: %v\n", id, err)
			} else {
				fmt.Printf("ID: %d, Response: %s\n", id, string(respBody))
			}

			resp.Body.Close()
		}(i)
	}

	wg.Wait()
	fmt.Println("Time Elapsed:", time.Since(start))

	fmt.Println("--------------------------------------------------------")
}
