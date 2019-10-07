package main

import (
	"log"
	"encoding/json"
	"os"
	"bufio"
	"io/ioutil"
	"strings"
)


// LOADING CONFIGURATION

// cannot adjust the JSON field names to be interoperable with https://arxiv.org/abs/1902.05260
type NodeInfo struct {
	Id			int			`json:"nid"`
	Ip			string		`json:"ip"`
	Port		int			`json:"port"`
	Capacity	float		`json:"cap"`
}

type TfReq struct {
	Src			int
	Dst			int
	Amount		float
}


func load_nodeinfo(filename string) NodeInfo {
	var nodeinfo NodeInfo

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Panicln("Error reading nodeinfo")
	}

	err = json.Unmarshal(data, &nodeinfo)
	if err != nil {
		log.Panicln("Error parsing nodeinfo")
	}

	return nodeinfo
}

func load_peers(filename string) []NodeInfo {
	var peers []NodeInfo

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Panicln("Error reading peers")
	}

	err = json.Unmarshal(data, &peers)
	if err != nil {
		log.Panicln("Error parsing peers")
	}

	return peers
}

func load_tfreqs(filename string) []TfReq {
	var tfreqs []TfReq

	file, err := os.Open(filename)
	if err != nil {
		log.Panicln("Error", err, "while reading transaction requests from file", filename)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		linefrag := strings.Split(string(line), ",")

		var tfreq TfReq
		tfreq.Src = myparseint(linefrag[0])
		tfreq.Dst = myparseint(linefrag[1])
		tfreq.Amount = myparsefloat(linefrag[2])

		tfreqs = append(tfreqs, tfreq)
	}

	if err := scanner.Err(); err != nil {
		log.Panicln("Error", err, "while parsing transaction requests")
	}

	return tfreqs
}

func load_paths(filename string) map[int][]Path {
	paths := make(map[int][]Path)

	file, err := os.Open(filename)
	if err != nil {
		log.Panicln("Error", err, "while reading paths from file", filename)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		linefrag := strings.Split(string(line), ",")

		var path Path
		dst_node_id := myparseint(linefrag[0])
		for _, v := range linefrag[1:] {
			path.Ids = append(path.Ids, myparseint(v))
		}

		paths[dst_node_id] = append(paths[dst_node_id], path)
	}

	if err := scanner.Err(); err != nil {
		log.Panicln("Error", err, "while parsing paths")
	}

	return paths
}
