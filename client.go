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

	pps := (bandWidth * 1000 * 1000) / 8000
	tick := 1000000 / pps
	fmt.Printf("band width %dmbps pps:%d tick %d micro\n", bandWidth, pps, tick)

	ticker := time.NewTicker(time.Duration(tick) * time.Microsecond)
	defer ticker.Stop()

	for range ticker.C {
		if sendStop {
			break
		}
		sendTs := time.Now().UnixNano()
		binary.BigEndian.PutUint64(data, uint64(sendTs))
		clientConn.Write(data)
		mCounter.addSendTs(uint64(sendTs))
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
		sendTs := binary.BigEndian.Uint64(recBuf)
		mCounter.addRecTs(sendTs)
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
	fmt.Printf("send:%d rec:%d loss rate:%f average rtt:%s\n", mCounter.getSendCount(), mCounter.getRecCount(), mCounter.lossRate(), mCounter.rtt())
}

type counter struct {
	sync.Mutex
	seqMap    map[uint64]uint64
	sendCount uint64
	recCount  uint64
}

func newCounter() *counter {
	c := new(counter)
	c.seqMap = make(map[uint64]uint64)
	c.sendCount = 0
	c.recCount = 0
	return c
}

func (c *counter) getRecCount() uint64 {
	c.Lock()
	defer c.Unlock()
	return c.recCount
}

func (c *counter) getSendCount() uint64 {
	c.Lock()
	defer c.Unlock()
	return c.sendCount
}

func (c *counter) addSendTs(ts uint64) {
	c.Lock()
	defer c.Unlock()

	c.sendCount = c.sendCount + 1
	c.seqMap[ts] = 0
}

func (c *counter) addRecTs(sendTs uint64) {
	c.Lock()
	defer c.Unlock()

	if _, ok := c.seqMap[sendTs]; ok {
		recTs := time.Now().UnixNano()
		c.seqMap[sendTs] = uint64(recTs)
	} else {
		fmt.Println("cant find ", sendTs)
	}
	c.recCount = c.recCount + 1
}

func (c *counter) lossRate() float32 {
	c.Lock()
	defer c.Unlock()

	var leftSize uint64
	for _, recTs := range c.seqMap {
		if recTs == 0 {
			leftSize = leftSize + 1
		}
	}

	if leftSize == 0 {
		return 0.0
	}

	lossRate := float32(leftSize) / float32(c.sendCount)
	return lossRate * 100
}

func (c *counter) rtt() string {
	c.Lock()
	defer c.Unlock()

	count := 0
	sumn := uint64(0)
	for sendTs, recTs := range c.seqMap {
		if recTs != 0 {
			rtt := recTs - sendTs
			count++
			sumn = sumn + rtt
		}
	}
	rttNano := float32(sumn) / float32(count)
	rttMill := rttNano / float32(1000000)
	return fmt.Sprintf("%f", rttMill)
}
