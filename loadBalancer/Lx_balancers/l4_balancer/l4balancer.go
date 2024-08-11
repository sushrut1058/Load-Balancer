package l4

import (
	"fmt"
	"io"
	"loadBalancer/store"
	"net"
)

func L4_balancer(conn net.Conn) {
	defer conn.Close()
	curUrl := store.GetUrl(store.Data["strategy"].(string))
	fmt.Println("Current url", curUrl)
	serverConn, err := net.Dial(store.Data["proto"].(string), curUrl.String())
	if err != nil {
		fmt.Printf("[l4balancer] Failed to connect to backend %s: %v", curUrl.String(), err)
		conn.Close()
		return
	}
	fmt.Println("ServerConn", serverConn)
	defer serverConn.Close()

	go io.Copy(serverConn, conn)
	io.Copy(conn, serverConn)
}
