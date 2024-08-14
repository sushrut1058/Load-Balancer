package main

import (
	"fmt"
	"io"
	"net/http"
	"sync"
)

func main() {
	fmt.Println("Started!")
	// file, err := os.ReadFile("body.txt")
	// if err != nil {
	// 	fmt.Println("Error reading file")
	// }

	var wg sync.WaitGroup

	client := &http.Client{
		Transport: &http.Transport{
			DisableKeepAlives: true, //crucial for testing, each request to be considered a different client.
		},
	}
	request, err := http.NewRequest("GET", "http://localhost:8080", nil)
	if err != nil {
		fmt.Println("", request, err)
	}
	k := 10
	for i := 0; i < k; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			resp, err := client.Do(request)
			if err != nil {
				fmt.Printf("Error in request: %v\n", err)
				return
			}
			fmt.Printf("Request sent. Status Code: %d\n", resp.StatusCode)
			body, error := io.ReadAll(resp.Body)
			if error != nil {
				fmt.Printf("Error while reading: %v\n", error)
			} else {
				// resp_string := string(body)
				fmt.Printf("[RESPONSE]: %v\n", body)
			}

			// io.Copy(io.Discard, resp.Body)
			resp.Body.Close()

		}()
	}

	wg.Wait()
	fmt.Println("--------------------------------------------------------")
}
