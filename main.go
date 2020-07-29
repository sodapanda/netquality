package main

import (
	"flag"
	"fmt"
	"os"
)

var isServer bool
var isClient bool
var bandWidth int

func main() {
	isServerFlag := flag.Bool("s", false, "server")
	isClientFlag := flag.Bool("c", false, "client")
	bandWidthFlag := flag.Int("b", 100, "bandwidth")

	flag.Parse()

	isServer = *isServerFlag
	isClient = *isClientFlag
	bandWidth = *bandWidthFlag

	if isServer {
		go startServer()
	}

	if isClient {
		go startClient()
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
