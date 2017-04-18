package main

import (
	"RMS_Node/Common"
	"RMS_Node/Serial_Srv"
	"RMS_Node/Xmodem"
	"fmt"
	"math/rand"
	"runtime"
	"time"
)

var RMSNode_EXIT chan int

func main() {
	runtime.GOMAXPROCS(2)
	Common.Init()

	//util.HRBserive(true)
	//util.HRBserive(false)

	//defer util.HRBserive(true)
	var dd []byte = make([]byte, 8)
	for i, _ := range dd {
		dd[i] = byte(rand.Int31n(255))
	}
	if len(dd) <= 1024 {
		fmt.Println(dd)
	} else {
		fmt.Println(dd[:1024])
	}
	go Xmodem.XmodemTransmit(dd)
	go Serial_Srv.SerialPortDaemon()
	for {
		time.Sleep(10e9)
	}
	<-RMSNode_EXIT
}
