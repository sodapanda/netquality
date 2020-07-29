package main

import (
	"container/list"
	"encoding/binary"
	"fmt"
	"net"
	"sync"
	"time"
)

var mCounter *counter

func startClient() {
	fmt.Println("client started")
	mCounter = newCounter()

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
	seqNum = 1
	for {
		binary.BigEndian.PutUint64(data, seqNum)
		clientConn.Write(data)
		mCounter.addSendSeq(seqNum)
		time.Sleep(waitTime * time.Millisecond)
		seqNum++
	}
}

func clientRecData(conn *net.UDPConn) {
	recBuf := make([]byte, 1400)
	for {
		_, err := conn.Read(recBuf)
		checkErr(err)
		seqNum := binary.BigEndian.Uint64(recBuf)
		mCounter.removeRecSeq(seqNum)
		mCounter.incRecCount()
	}
}

func printLog() {
	for {
		time.Sleep(1 * time.Second)
		fmt.Printf("send:%d rec:%d loss rate:%f\n", mCounter.getSendCount(), mCounter.getRecCount(), mCounter.lossRate())
	}
}

type counter struct {
	sync.Mutex
	seqList    *list.List
	maxRecSeq  uint64
	maxSendSeq uint64
	recCount   uint64
}

func newCounter() *counter {
	c := new(counter)
	c.seqList = list.New()
	c.maxRecSeq = 0
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
	c.seqList.PushBack(seq)
}

func (c *counter) removeRecSeq(seq uint64) {
	c.Lock()
	defer c.Unlock()

	if seq > c.maxRecSeq {
		c.maxRecSeq = seq
	}

	var toRemove *list.Element
	toRemove = nil

	for e := c.seqList.Front(); e != nil; e = e.Next() {
		if e.Value == seq {
			toRemove = e
			break
		}
	}

	if toRemove != nil {
		c.seqList.Remove(toRemove)
	}
}

func (c *counter) lossRate() float32 {
	c.Lock()
	defer c.Unlock()

	leftSize := uint64(c.seqList.Len())
	if leftSize == 0 {
		return 0.0
	}

	inFlightSize := c.maxSendSeq - c.maxRecSeq

	lossRate := float32(leftSize-inFlightSize) / float32(c.maxRecSeq)

	c.seqList = list.New()
	return lossRate * 100
}
