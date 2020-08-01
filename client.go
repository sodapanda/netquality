package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"sync"
	"time"
)

var mCounter *counter
var sendStop bool
var recStop bool

func startClient() {
	fmt.Println("client started")
	mCounter = newCounter()
	sendStop = false
	recStop = false

	serverAddr, err := net.ResolveUDPAddr("udp4", serverIP+":"+serverPort)
	checkErr(err)

	clientConn, err := net.DialUDP("udp4", nil, serverAddr)
	checkErr(err)
	// defer clientConn.Close()

	go clientRecData(clientConn)
	go stopWatch()

	data := make([]byte, 10)
	for i := range data {
		data[i] = 1
	}

	var seqNum uint64
	seqNum = 1

	pps := (bandWidth * 1000 * 1000) / 8000
	tick := 1000000 / pps
	fmt.Printf("band width %dmbps pps:%d tick %d micro\n", bandWidth, pps, tick)

	ticker := time.NewTicker(time.Duration(tick) * time.Microsecond)
	defer ticker.Stop()

	for range ticker.C {
		if sendStop {
			break
		}
		binary.BigEndian.PutUint64(data, seqNum)
		clientConn.Write(data)
		mCounter.addSendSeq(seqNum)
		seqNum++
	}
}

var largestPacketSize int

func clientRecData(conn *net.UDPConn) {
	recBuf := make([]byte, 1400)
	largestPacketSize = 0
	for {
		if recStop {
			break
		}
		length, err := conn.Read(recBuf)
		if length > largestPacketSize {
			largestPacketSize = length
		}
		checkErr(err)
		seqNum := binary.BigEndian.Uint64(recBuf)
		mCounter.removeRecSeq(seqNum)
		mCounter.incRecCount()
	}
}

func stopWatch() {
	time.Sleep(10 * time.Second)
	sendStop = true
	time.Sleep(1 * time.Second)
	recStop = true
	time.Sleep(100 * time.Millisecond)
	printLog()
}

func printLog() {
	fmt.Printf("send:%d rec:%d loss rate:%f\n", mCounter.getSendCount(), mCounter.getRecCount(), mCounter.lossRate())
}

type counter struct {
	sync.Mutex
	seqMap     map[uint64]bool
	maxSendSeq uint64
	recCount   uint64
}

func newCounter() *counter {
	c := new(counter)
	c.seqMap = make(map[uint64]bool)
	c.maxSendSeq = 0
	return c
}

func (c *counter) incRecCount() {
	c.Lock()
	defer c.Unlock()

	c.recCount = c.recCount + 1
}

func (c *counter) getRecCount() uint64 {
	c.Lock()
	defer c.Unlock()
	return c.recCount
}

func (c *counter) getSendCount() uint64 {
	c.Lock()
	defer c.Unlock()
	return c.maxSendSeq
}

func (c *counter) addSendSeq(seq uint64) {
	c.Lock()
	defer c.Unlock()

	c.maxSendSeq = seq
	c.seqMap[seq] = false
}

func (c *counter) removeRecSeq(seq uint64) {
	c.Lock()
	defer c.Unlock()

	if _, ok := c.seqMap[seq]; ok {
		c.seqMap[seq] = true
	} else {
		fmt.Println("cant find ", seq)
	}
}

func (c *counter) lossRate() float32 {
	c.Lock()
	defer c.Unlock()

	var leftSize uint64
	for _, v := range c.seqMap {
		if !v {
			leftSize = leftSize + 1
		}
	}

	if leftSize == 0 {
		return 0.0
	}

	lossRate := float32(leftSize) / float32(c.maxSendSeq)
	return lossRate * 100
}
