package main

import (
	"log"
	"time"
	"net"
	"bufio"
	"encoding/json"
)


// shovel messages around

func handle_inbound_connections(conn *net.TCPListener, queue_routing chan P2pMsg) {
	for {
		conn_client, err := conn.AcceptTCP()
		if err != nil {
			log.Fatalln("Error accepting connection:", err)
		}

		conn_client.SetKeepAlive(true)
		conn_client.SetKeepAlivePeriod(DELAY_KEEPALIVE_S * time.Second)
		conn_client.SetNoDelay(true)

		go handle_inbound_session(conn_client, queue_routing)
	}
}

func handle_inbound_session(conn *net.TCPConn, queue_routing chan P2pMsg) {
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		line := scanner.Text()
		// log.Println("Received line on inbound session:", line)

		var msg P2pMsg
		err := json.Unmarshal(([]byte)(line), &msg)
		if err != nil {
			log.Panicln("Error parsing inbound data")
		}

		go func() {
			// simulate delay for fast/slow packets
			if is_fast(msg) {
				// log.Println("Fast packet!")
				time.Sleep(time.Duration(myrandrangei(DELAY_NETWORK_FAST_MIN_MS, DELAY_NETWORK_FAST_MAX_MS)) * time.Millisecond)
			} else {
				// log.Println("Slow packet!")
				time.Sleep(time.Duration(myrandrangei(DELAY_NETWORK_SLOW_MIN_MS, DELAY_NETWORK_SLOW_MAX_MS)) * time.Millisecond)
			}
			queue_routing <- msg
		} ()
	}

	if err := scanner.Err(); err != nil {
		log.Panicln("Error", err, "while handling inbound session")
	}
}

func handle_outbound_session(session *Session) {
	for msg := range session.QueueTransmit {
		// log.Println("Outbound msg:", msg)

		data, err := json.Marshal(msg)
		if err != nil {
			log.Panicln("Error serializing outbound data", err)
		}

		write_len, err := session.Connection.Write(append(data, ([]byte)("\n")...))
		if err != nil || write_len != len(data)+1 {
			session.Connection.Close()
			log.Panicln("Error transmitting outbound data", err, write_len, len(data))
		}
	}
}
