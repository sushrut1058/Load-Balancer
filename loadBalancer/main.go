package main

import (
	"encoding/json"
	"fmt"
	l4 "loadBalancer/Lx_balancers/l4_balancer"
	l7 "loadBalancer/Lx_balancers/l7_balancer"
	"loadBalancer/caching/basicCaching"
	rCaching "loadBalancer/caching/redis"
	store "loadBalancer/store"
	types "loadBalancer/types"
	"net"
	"net/http"
	"os"
)

func readConfig() {
	file, err := os.ReadFile("conf.json")
	if err != nil {
		fmt.Printf("[main] Error while reading json: %v\n", err)
	}
	err = json.Unmarshal(file, &store.Data)
	if err != nil {
		fmt.Printf("[main] Error while writing json to memory: %v\n", err)
	}

	//caching
	if store.Data["caching"] == "baseline" {
		store.SetCache = basicCaching.SetCache
		store.GetCachedResponse = basicCaching.GetCachedResponse
	} else if store.Data["caching"] == "redis" {
		store.SetCache = rCaching.SetCache
		store.GetCachedResponse = rCaching.GetCachedResponse
	} else {
		store.SetCache = func(string, []byte, *http.Response) {}
		store.GetCachedResponse = func(string) (*types.Cache, bool) { return nil, false }
	}

	//servers and port
	store.Servers, _ = store.Data["servers"].([]interface{})
	store.Port = fmt.Sprintf(":%v", store.Data["port"].(string))

}

func main() {
	fmt.Println("# Starting Application...")
	readConfig()
	fmt.Println("# Configuration stored successfully.")
	if store.Data["level"] == "L7" {
		fmt.Println("[main l7] starting")
		http.HandleFunc("/", l7.L7_balancer)
		http.ListenAndServe(store.Port, nil)
		fmt.Println("[main l7] listening")
	} else {
		listener, err := net.Listen(store.Data["proto"].(string), store.Port)
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
			go l4.L4_balancer(conn)
		}
	}
}
