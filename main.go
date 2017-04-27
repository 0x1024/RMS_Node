package main

import (
	"RMS_Node/Common"
	"RMS_Node/Serial_Srv"
	"RMS_Srv/ExtPortSrv"
	"runtime"
	"time"
)

var RMSNode_EXIT chan int
var RMSNode_EXIT1 chan int

func main() {
	runtime.GOMAXPROCS(2)
	Common.Init()

	//util.HRBserive(true)
	go ExtPortSrv.NodeStarter()

	go Serial_Srv.SerialPortDaemon()

	go ExtPortSrv.NodeStarter()


	for {
		time.Sleep(10e9)
	}

	<-RMSNode_EXIT
}

