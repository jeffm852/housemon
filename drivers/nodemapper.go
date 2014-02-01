package drivers

import (
	"fmt"
	"encoding/json"
	"os"
)

type NodeMapEntry struct {
	G      int
	N      int
	Ts     int64
	Te     int64
	Type   string
	Name   string
}

type NodeMapT struct {
	NodeMap []NodeMapEntry
}

var nMap *NodeMapT

func jNodeMap() {
	file, e := os.Open("./nodemap.json")
	check(e)
	decoder := json.NewDecoder(file)
	decoder.Decode(&nMap)
	//fmt.Printf("Results: %v\n", nMap)
}

func JNodeType(B, G, N int, TS int64) (bool, string, string) {
	for _, node := range nMap.NodeMap {
		//TODO: Add band in identification check
		B += 0
		if node.G == G && node.N == N && TS >= node.Ts && (TS < node.Te || node.Te <= node.Ts) {
			fmt.Println("found rf12 node: ", node.Type)
			return true, node.Type, node.Name
		} else {
		}
	}
	return false, "", ""
}

func init() {
	jNodeMap()
}
