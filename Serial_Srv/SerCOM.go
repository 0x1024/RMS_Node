package Serial_Srv

import (
	"RMS_Node/util"
	"fmt"
	"github.com/juju/errors"
	log "github.com/sirupsen/logrus"
	"go.bug.st/serial.v1"
)

var port serial.Port
var ComPortDown_ch chan int
var ctr uint

// Open the first serial port detected at 9600bps N81
var mode = &serial.Mode{
	BaudRate: 115200,
	Parity:   serial.NoParity,
	DataBits: 8,
	StopBits: serial.OneStopBit,
}

func SerialPortDaemon() {
	for port == nil {
		//enum ports
		ports, err := serial.GetPortsList()
		if err != nil {
			log.Fatal(err)
		}
		if len(ports) == 0 {
			log.Fatal("No serial ports found!")
		}
		// Print the list of detected ports
		for _, port := range ports {
			fmt.Printf("Found port: %v\n", port)
		}

		for _, portx := range ports {

			port, err = serial.Open(portx, mode)
			if err != nil {
				log.Fatal(err)
			}
			for {
				echo(port)
				EchoWaiter(port)
			}
		}

		fmt.Println("\nserv end ")
		<-ComPortDown_ch
	}

}

func echo(port serial.Port) {

	send := make([]byte, 32)
	send[0] = 0xaa
	send[1] = 0x55
	send[2] = 0x01
	send[3] = 0x00
	send[4] = 0x00
	send[5] = 0xF0
	send[6] = 0x80
	ret := util.CRC16(send, 7)
	send[7] = byte((ret >> 8) & 0xff)
	send[8] = byte(ret & 0xff)

	n, err := port.Write(send[:9])
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Handshake echo, %d bytes:%x \n", n, send[:9])

}

func EchoWaiter(port serial.Port) { // Read and print the response

	buff := make([]byte, 4096)
	var res []byte
	var cnt int = 0

Out:
	for true {
		for {
			// Reads up to 100 bytes
			n, err := port.Read(buff)

			if err != nil {
				log.Fatal(err)
				break
			}
			if n == 0 {

				log.Println("\nEOF")
				break
			}

			cnt++
			//for _, as := range buff[:n] {
			//	fmt.Printf("%02X ", as)
			//}
			res = append(res, buff[:n]...)
			//fmt.Printf("\r\nrec %d     \r\n", cnt, res)
			if len(res) > 7 {
				err := DeFrame(res)
				if err == nil {
					break Out
				} else if err.Error() =="1" {

				}else {
					fmt.Println(err)
					panic(err)
				}

			}

		}

	}
}

func DeFrame(res []byte) error {
	if res[0] == 0xAA && res[1] == 0x55 && len(res) == int(9+res[4]) {
		ctr++
		fmt.Println("\ngot pack",ctr, res)
		return nil
	}else if res[0] == 0xAA && res[1] == 0x55 && len(res) < int(9+res[4]){
		fmt.Println("\nshort pack",ctr, res)
		return errors.New("1")
	}
	fmt.Println("\nwrong pack", res)
	return errors.New("wrong Frame~!")
}
