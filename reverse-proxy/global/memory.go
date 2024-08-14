package global

var Data map[string]interface{}
var RequestQueue Queue
var NServers int
var Servers []Resource
var ServerIndexMap map[string]int
var CurrentCapacity []int
var UrlIndex uint32
var MaxWorkerCount int
