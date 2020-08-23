package main

import (
	"fmt"
	"testing"

	"github.com/klauspost/reedsolomon"
)

func TestFec(t *testing.T) {
	enc, err := reedsolomon.New(4, 2)
	if err != nil {
		fmt.Println(err)
	}
	allData := make([][]byte, 6)
	for i := 0; i < 6; i++ {
		allData[i] = make([]byte, 10)
		for j := 0; j < 10; j++ {
			allData[i][j] = byte(i + j)
		}
	}
	err = enc.Encode(allData)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(allData)
	}

	allData[0] = nil
	allData[4] = nil
	err = enc.Reconstruct(allData)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(allData)
	}
}
