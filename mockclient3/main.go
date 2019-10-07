package main

import (
	"log"
	"time"
	"net"
	"os"
	"github.com/davecgh/go-spew/spew"
)


// FILTER TRANSACTION REQUESTS

func filter_tfreqs_not(tfreqs_filtered chan TfReq, tfreqs []TfReq) {
	for _, tfreq := range tfreqs {
		tfreqs_filtered <- tfreq
	}

	close(tfreqs_filtered)
}

func filter_tfreqs_initially_locally_feasible(tfreqs_filtered chan TfReq, tfreqs []TfReq, peers []NodeInfo) {
	var total_capacity float

	for _, p := range peers {
		total_capacity += p.Capacity
	}

	for _, tfreq := range tfreqs {
		if tfreq.Amount <= total_capacity {
			tfreqs_filtered <- tfreq
		}
	}

	close(tfreqs_filtered)
}


// COMPILED CONFIG

const NUM_QUEUE_SIZE = 1024*32
const NUM_BUFFER_SIZE = 16*2014

const DELAY_AFTER_LISTENING_S = 5
const DELAY_AFTER_CONNECTING_S = 5
const DELAY_KEEPALIVE_S = 5

const DELAY_NETWORK_SLOW_MIN_MS = 50
const DELAY_NETWORK_SLOW_MAX_MS = 150
const DELAY_NETWORK_FAST_MIN_MS = 0
const DELAY_NETWORK_FAST_MAX_MS = 0


// MAIN PROGRAM

func main() {

	// setup logging
	log.SetOutput(os.Stdout)


	// setup queues
	queue_routing := make(chan P2pMsg, NUM_QUEUE_SIZE)
	queue_local := make(chan P2pMsg, NUM_QUEUE_SIZE)


	// parse commandline arguments and load config

	if len(os.Args) < 6+1 {
		log.Fatalln("Usage: mockclient <${n}.json> <n${n}.json> <graph file> <tr${n}.txt> <pa${n}.txt> <algo> [algorithm specific parameters]")
		// 1 = node config
		// 2 = peers
		// 3 = graph
		// 4 = transactions
		// 5 = paths
		// 6 = algorithm
	}

	nodeinfo := load_nodeinfo(os.Args[1])
	peers := load_peers(os.Args[2])
	// graph is not loaded at this point
	tfreqs := load_tfreqs(os.Args[4])
	paths := load_paths(os.Args[5])

	algo := os.Args[6]
	var algo_param_paths int
	var algo_param_retries int
	var algo_param_redundancy int
	if algo == "retry-02-amp-2" {
		if len(os.Args) != 6+2+1 {
			log.Fatalln("Parameters for 'retry-02-amp-2' algorithm: <number of paths> <number of retries>")
		} else {
			algo_param_paths = myparseint(os.Args[7])
			algo_param_retries = myparseint(os.Args[8])
		}
	} else if algo == "redundancy-02-amp-2" {
		if len(os.Args) != 6+2+1 {
			log.Fatalln("Parameters for 'redundancy-02-amp-2' algorithm: <number of paths> <number of redundant transactions>")
		} else {
			algo_param_paths = myparseint(os.Args[7])
			algo_param_redundancy = myparseint(os.Args[8])
		}
	} else if algo == "redundantretry-02-amp-2" {
		if len(os.Args) != 6+3+1 {
			log.Fatalln("Parameters for 'redundantretry-02-amp-2' algorithm: <number of paths> <number of retries> <number of redundant transactions>")
		} else {
			algo_param_paths = myparseint(os.Args[7])
			algo_param_retries = myparseint(os.Args[8])
			algo_param_redundancy = myparseint(os.Args[9])
		}
	} else {
		log.Println("Available algorithms:")
		log.Println("  retry-02-amp-2           ...")
		log.Println("  redundancy-02-amp-2      ...")
		log.Println("  redundantretry-02-amp-2  ...")
		log.Fatalln("Unknown algorithm:", algo)
	}

	log.Println("Configuration loaded")
	spew.Dump(nodeinfo)
	spew.Dump(peers)
	// spew.Dump(tfreqs[0:5])
	spew.Dump(paths)
	spew.Dump(algo)


	// listen for inbound connections

	socket_tcp_inbound, err := net.ListenTCP("tcp4", &net.TCPAddr{net.ParseIP(nodeinfo.Ip), nodeinfo.Port, ""})
	if err != nil {
		log.Fatalln("Failed to open port:", err)
	}
	defer socket_tcp_inbound.Close()
	go handle_inbound_connections(socket_tcp_inbound, queue_routing)
	log.Println("Listening to network")


	// wait for all clients to start listening before connecting to the peers

	log.Println("Waiting for other clients to finish startup ...")
	time.Sleep(DELAY_AFTER_LISTENING_S * time.Second)
	log.Println("Done!")


	// connect to peers

	var sessions_outbound []Session

	for _, peer := range peers {
		socket_tcp_outbound, err := net.DialTCP("tcp4", nil, &net.TCPAddr{net.ParseIP(peer.Ip), peer.Port, ""})
		if err != nil {
			log.Fatalln("Error ", err, "connecting to peer", peer.Ip, peer.Port)
		}

		socket_tcp_outbound.SetKeepAlive(true)
		socket_tcp_outbound.SetKeepAlivePeriod(DELAY_KEEPALIVE_S*time.Second)
		socket_tcp_outbound.SetNoDelay(true)

		var s Session
		s.Peer = peer
		s.QueueTransmit = make(chan P2pMsg, NUM_QUEUE_SIZE)
		s.Connection = socket_tcp_outbound
		s.Capacity = peer.Capacity

		sessions_outbound = append(sessions_outbound, s)
	}

	for i, s := range sessions_outbound {
		log.Println(i)
		spew.Dump(s)
		go handle_outbound_session(&(sessions_outbound[i]))
	}


	// start local routing

	go handle_routing(nodeinfo, queue_routing, queue_local, sessions_outbound)


	// wait for all clients to connect to peers before processing transactions

	log.Println("Waiting for other clients to finish connecting ...")
	time.Sleep(DELAY_AFTER_CONNECTING_S * time.Second)
	log.Println("Done!")


	// filter transactions

	tfreqs_filtered := make(chan TfReq, 1)
	go filter_tfreqs_not(tfreqs_filtered, tfreqs)
	// go filter_tfreqs_initially_locally_feasible(tfreqs_filtered, tfreqs, peers)


	// process transactions

	if algo == "retry-02-amp-2" {
		process_tfreqs_retry_02_amp_2(nodeinfo, tfreqs_filtered, paths, peers, sessions_outbound, queue_routing, queue_local, algo_param_paths, algo_param_retries)
	} else if algo == "redundancy-02-amp-2" {
		process_tfreqs_redundancy_02_amp_2(nodeinfo, tfreqs_filtered, paths, peers, sessions_outbound, queue_routing, queue_local, algo_param_paths, algo_param_redundancy)
	} else if algo == "redundantretry-02-amp-2" {
		process_tfreqs_redundantretry_02_amp_2(nodeinfo, tfreqs_filtered, paths, peers, sessions_outbound, queue_routing, queue_local, algo_param_paths, algo_param_retries, algo_param_redundancy)
	} else {
		log.Fatalln("Unknown algorithm:", algo)
	}


	// done, infinite loop to provide routing to the network

	select {}

}
