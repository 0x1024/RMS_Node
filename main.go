package main

import (
	"RMS_Node/Serial_Srv"

	"RMS_Node/util"
)
var RMSNode_EXIT chan int

func main() {
	//util.HRBserive(true)
	util.HRBserive(false)
	go Serial_Srv.SerialPortDaemon()
	defer util.HRBserive(true)

	<- RMSNode_EXIT
}







