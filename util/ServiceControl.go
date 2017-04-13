package util

import (
	"fmt"
	"os/exec"
)

func HRBserive(SetRun bool) {
	var cmd *exec.Cmd

	if SetRun {
		cmd = exec.Command("net", "start", "HitRobotBase")
	} else {
		cmd = exec.Command("net", "stop", "HitRobotBase")

	}

	err := cmd.Run()
	if err !=nil{
		fmt.Printf("%+v",err)
		//panic("\ncmd run with fault\n\n")
	}

	cmd.Wait()


}
