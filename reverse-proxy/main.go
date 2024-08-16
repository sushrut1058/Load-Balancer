package main

import (
	"fmt"
	"net"
	"net/http"
	"reverse-proxy/caching"
	global "reverse-proxy/global"
	"reverse-proxy/requests"
	"reverse-proxy/worker"
)

func handler_l7(w http.ResponseWriter, r *http.Request) {
	fmt.Println("request reached")
	done := make(chan bool)
	newRequestHandle := &requests.HTTPRequestHandle{Request: r, Writer: w, Processed: &done}
	requests.RequestChannel <- newRequestHandle
	<-done
}

func handler_l4(conn net.Conn) {
	newRequestHandle := &requests.TCPRequestHandle{Conn: conn}
	requests.RequestChannel <- newRequestHandle
}

func main() {
	fmt.Println("Starting . . .")

	global.Preprocessing()

	caching.InitCaching()

	worker.StartWorkerPool()

	if global.Data["level"] == "L7" {
		http.HandleFunc("/", handler_l7)
		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			fmt.Println("Error while listening. error:", err)
		}
	} else if global.Data["level"] == "L4" {
		listener, err := net.Listen(global.Data["proto"].(string), ":8080")
		fmt.Println("Test", global.Data["proto"].(string), global.Data["port"].(string))
		if err != nil {
			fmt.Printf("[main l4] Error creating listener: %v\n", err)
		} else {
			fmt.Println("[main l4] Listener created successfully")
		}
		defer listener.Close()
		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Printf("[main l4] Error accepting connection: %v\n", err)
				continue
			}
			go handler_l4(conn)
		}
	} else {
		fmt.Println("Error! Please check the loadbalancer type")
	}

}
