package requests

import "reverse-proxy/global"

var LoadBalancingStrategy func(index uint32) uint32

func InitLoadBalancing() {
	if global.Data["strategy"] == "round-robin" {
		LoadBalancingStrategy = roundRobin
	} else {
		LoadBalancingStrategy = weightedRoundRobin
	}
}

func weightedRoundRobin(index uint32) uint32 {
	if global.TotalCapacity[index] == counter {
		counter = 0
		UrlIndex++
	}
	index = UrlIndex % uint32(global.NServers)
	counter++
	return index
}

func roundRobin(index uint32) uint32 {
	UrlIndex++
	return (index + 1) % uint32(global.NServers)
}
