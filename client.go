package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"time"
)

func startClient() {
	fmt.Println("client started")

	serverAddr, err := net.ResolveUDPAddr("udp4", "127.0.0.1:21007")
	checkErr(err)

	clientConn, err := net.DialUDP("udp4", nil, serverAddr)
	checkErr(err)
	defer clientConn.Close()

	go clientRecData(clientConn)

	data := make([]byte, 10)
	for i := range data {
		data[i] = 1
	}

	var seqNum uint64
	seqNum = 0
	for {
		binary.BigEndian.PutUint64(data, seqNum)
		clientConn.Write(data)
		time.Sleep(waitTime * time.Millisecond)
		seqNum++
	}
}

func clientRecData(conn *net.UDPConn) {
	recBuf := make([]byte, 1400)
	for {
		_, err := conn.Read(recBuf)
		checkErr(err)
		// seqNum := binary.BigEndian.Uint64(recBuf)
	}
}
