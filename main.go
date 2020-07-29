package main

import (
	"flag"
	"fmt"
	"os"
	"time"
)

var isServer bool
var isClient bool
var waitTime time.Duration

func main() {
	isServerFlag := flag.Bool("s", false, "server")
	isClientFlag := flag.Bool("c", false, "client")
	waitTimeFlag := flag.Int("w", 100, "wait time")

	flag.Parse()

	isServer = *isServerFlag
	isClient = *isClientFlag
	waitTime = time.Duration(*waitTimeFlag)

	if isServer {
		go startServer()
	}

	if isClient {
		go startClient()
		go printLog()
	}

	var input string
	fmt.Scanln(&input)
}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
