package main

import (
	"fmt"
	"net"
)

func startServer() {
	fmt.Println("server started")

	serverAddr, err := net.ResolveUDPAddr("udp4", ":21007")
	checkErr(err)

	serverConn, err := net.ListenUDP("udp4", serverAddr)
	checkErr(err)
	defer serverConn.Close()

	recBuf := make([]byte, 100)
	buf := make([]byte, 1400)
	for i := range buf {
		buf[i] = 1
	}

	for {
		length, addr, err := serverConn.ReadFromUDP(recBuf)
		checkErr(err)
		data := recBuf[:length]
		//数据拿出来，发包序号 补上站位的数据，发送回去，带着序号
		copy(buf, data[:8])
		length, err = serverConn.WriteToUDP(buf, addr)
		checkErr(err)
	}
}
