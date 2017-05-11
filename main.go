package main

import (
	"RMS_Node/Common"
	"RMS_Node/Serial_Srv"
	"RMS_Node/util"
	"RMS_Srv/ExtPortSrv"
	"RMS_Srv/Public"
	"fmt"
	"runtime"
)

var RMSNode_EXIT chan int
var RMSNode_EXIT1 chan int

var buildstamp = "no timestamp set"
var githash = "no githash set"
var project = "RMS_Node"

func main() {
	fmt.Println(project)
	fmt.Println(" Buildstamp is:", buildstamp)
	fmt.Println("Buildgithash is:", githash)

	runtime.GOMAXPROCS(2)
	Common.Init()
	Public.Init()

	util.HRBserive(false)
	go ExtPortSrv.NodeStarter()

	go Serial_Srv.SerialPortDaemon()

	<-RMSNode_EXIT
}
