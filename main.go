package main

import (
	"flag"
	"fmt"
	"os"
)

var isServer bool
var isClient bool
var bandWidth int
var serverIP string
var serverPort string

func main() {
	isServerFlag := flag.Bool("s", false, "server")
	isClientFlag := flag.Bool("c", false, "client")
	bandWidthFlag := flag.Int("b", 100, "bandwidth")
	serverIPFlag := flag.String("ip", "127.0.0.1", "remote ip")
	serverPortFlag := flag.String("p", "21007", "server port")

	flag.Parse()

	isServer = *isServerFlag
	isClient = *isClientFlag
	bandWidth = *bandWidthFlag
	serverIP = *serverIPFlag
	serverPort = *serverPortFlag

	if isServer {
		go startServer()
	}

	if isClient {
		go startClient()
	}

	var input string
	fmt.Scanln(&input)

	if isClient {
		fmt.Println("max packet size ", largestPacketSize)
	}
}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
