package global

type Resource struct {
	URL      string
	Capacity float64
}

var Data map[string]interface{}
var NServers int
var Servers []Resource
var ServerIndexMap map[string]int
var CurrentCapacity []int
var TotalCapacity []int
var MaxWorkerCount int

var CurrentWorkerCount int32
