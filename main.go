package main

import (
	"RMS_Node/Common"
	"RMS_Node/Serial_Srv"
	//	"RMS_Srv/ExtPortSrv"
	"runtime"
	"time"
)

var RMSNode_EXIT chan int

func main() {
	runtime.GOMAXPROCS(2)
	Common.Init()

	//util.HRBserive(true)
	//	go ExtPortSrv.ExternService()

	go Serial_Srv.SerialPortDaemon()
	for {
		time.Sleep(10e9)
	}
	<-RMSNode_EXIT
}
