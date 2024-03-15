package main

import (
	"fmt"
	"net/http"
	"os"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("[PORT: %s] Request received on server\n", os.Args[1])
	resp := fmt.Sprintf("Hello from server at Port:%v", os.Args[1])
	fmt.Fprintf(w, resp)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("No port passed in args")
		return
	}
	http.HandleFunc("/", handler)
	port := os.Args[1]
	fmt.Printf("Starting on port: %s\n", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Printf("Error received in ListenAndServe: %s\n", err)
	}
}
