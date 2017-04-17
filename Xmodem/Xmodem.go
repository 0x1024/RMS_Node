package Xmodem

import(

)
import (
	"RMS_Node/util"
	"RMS_Node/Common"
	"RMS_Node/Serial_Srv"
	"fmt"
	"time"
)

const(
       SOH    =0x01
       STX    =0x02
       EOT    =0x04
       ACK    =0x06
       NAK    =0x15
       CAN    =0x18
       CTRLZ  =0x1A
)
const MAXRETRANS int =25
var Ch_GetByte chan byte
var Ch_ReqByte chan byte
var Ch_DumpAll chan byte
var RecPool []byte
func Inpool(){
	for {
		RecPool = append(RecPool, <-Common.Ch_ComStreamData...)
		fmt.Println("inpool ",RecPool)

	}
}
func  _inbyte() byte{
	var to int =0
	for len(RecPool) == 0 {
		time.Sleep(1e6)
		if to++;to > 100 {return 0}
	}
	if RecPool[0] == 6{
		fmt.Println("inb " , RecPool)
	}
	ret := RecPool[0]
	RecPool = RecPool[1:]
	fmt.Println("after wait",ret)
	return ret
}
func  _outbyte(rec byte) {
	var send []byte=make([]byte,1)
	send[0] = rec
	Serial_Srv.SendByte( send)

}
func _outSteam(send []byte){
	Serial_Srv.SendByte(send)
}

func flushinput(){
	RecPool=RecPool[len(RecPool):]
}
func XmodemTransmit(src []byte)int{
	go Inpool()
	rec:=<- Common.Ch_Ft_start
	fmt.Println("x rec",rec)
	var xbuff []byte=make([]byte,1030)
	var bufsz int
	var crc uint16
	var retry int
	var packetno byte=1
	var i, lens int =0,0
	var c int=0
	var srcsz int=len(src)
	//part 1
	for retry=0; retry < 16; retry++ {
		fmt.Println("xt handshake")
		if cc := _inbyte();cc>=0 {
		switch cc {
			case 'C':
				crc = 1
				flushinput()
				goto start_trans_head
			case NAK:
				crc = 0
				goto start_trans_head
			case CAN:
				if cc = _inbyte();cc == CAN {
				_outbyte(ACK)
				flushinput()
				return -1 /* canceled by remote */
			}
				break
			default:
				break
			}
		}
	}
	_outbyte(CAN)
	_outbyte(CAN)
	_outbyte(CAN)
	flushinput()
	fmt.Println("ft quit -2")
	return -2 /* no sync */




start_trans_head:
	xbuff[0] = STX
	bufsz = 1024
	xbuff[1] = 0
	xbuff[2] = 255 -0


	//file name,lens
	result := []byte("012345 ")
	result =append(result,0)
	result =append(result,[]byte("8")... )
	copy(xbuff[3:],result)

	ccrc := util.Crc16_ccitt(xbuff[3:3+bufsz],bufsz)
	xbuff[bufsz+3] = byte(ccrc>>8) & 0xFF
	xbuff[bufsz+4] = byte(ccrc & 0xFF)

	for retry = 0; retry < MAXRETRANS;retry++ {
		flushinput()
		_outSteam(xbuff[:bufsz+5])
		time.Sleep(5e4)

		for c = int( _inbyte() );c > 0;  c = int( _inbyte() ){

			switch c {
			case ACK:
				//lens += bufsz
				flushinput()

				goto start_trans
			case CAN:
				if c = int( _inbyte() ) ;c== CAN {
					_outbyte(ACK)
					flushinput()
					return -1 /* canceled by remote */
				}
				break
			case NAK:

			default:
				break
			}
			//time.Sleep(5e6)

		}
	}

start_trans:
	fmt.Println("\nstart trans\n\n\n")
	for {
		xbuff[0] = STX
		bufsz = 1024
		xbuff[1] = packetno
		xbuff[2] = 255 -packetno
		c = srcsz - lens
		if c > bufsz {	c = bufsz}

		if c >= 0 {
			if c == 0 {
				xbuff[3] = CTRLZ
			}else {
				copy (xbuff[3:], src[lens:lens+c])
				if c < bufsz{ xbuff[3+c] = CTRLZ}
			}

			if crc>0 {
				ccrc := util.Crc16_ccitt(xbuff[3:3+bufsz],bufsz)
				xbuff[bufsz+3] = byte(ccrc>>8) & 0xFF
				xbuff[bufsz+4] = byte(ccrc & 0xFF)
				fmt.Println("bc  ",xbuff[bufsz:bufsz+5])

			}else {
				ccks := 0
				for i = 3; i < bufsz+3; i++ {
					ccks = ccks+ int(xbuff[i])
				}
				xbuff[bufsz+3] = byte(ccks)
			}

			for retry = 0; retry < MAXRETRANS;retry++ {

				_outSteam(xbuff[:bufsz+5])
				fmt.Println("\nXm pack send",retry)
				time.Sleep(5e4)
				for c = int( _inbyte() );c > 0;  c = int( _inbyte() ){
					switch c {
						case ACK:
							packetno++
							lens += bufsz
							goto start_trans
						case CAN:
							if c = int( _inbyte() ) ;c== CAN {
							_outbyte(ACK)
							flushinput()
							return -1 /* canceled by remote */
							}
						break
						case NAK:

						default:
						break
					}
				}
			}
			_outbyte(CAN)
			_outbyte(CAN)
			_outbyte(CAN)
			flushinput()
			return -4 /* xmit error */





		}else {
			fmt.Println("EOT")
			for retry = 0; retry < 10; retry++ {
			_outbyte(EOT)
			if c = int( _inbyte() )<<1 ;c== ACK {break}
			}
			flushinput()
				if c == ACK {
					return lens
				}else {
					return -5
				}
		}

	}



return 0
}

func ret_bool(f int)int {
	if f > 0{
		return 1
	}
	return 0
}
//
//
//int xmodemReceive(unsigned char *dest, int destsz)
//{
//unsigned char xbuff[1030]; /* 1024 for XModem 1k + 3 head chars + 2 crc + nul */
//unsigned char *p;
//int bufsz, crc = 0;
//unsigned char trychar = 'C';
//unsigned char packetno = 1;
//int i, c, len = 0;
//int retry, retrans = MAXRETRANS;
//
//for(;;) {
//for( retry = 0; retry < 16; ++retry) {
//if (trychar) _outbyte(trychar);
//if ((c = _inbyte((DLY_1S)<<1)) >= 0) {
//switch (c) {
//case SOH:
//bufsz = 128;
//goto start_recv;
//case STX:
//bufsz = 1024;
//goto start_recv;
//case EOT:
//flushinput();
//_outbyte(ACK);
//return len; /* normal end */
//case CAN:
//if ((c = _inbyte(DLY_1S)) == CAN) {
//flushinput();
//_outbyte(ACK);
//return -1; /* canceled by remote */
//}
//break;
//default:
//break;
//}
//}
//}
//if (trychar == 'C') { trychar = NAK; continue; }
//flushinput();
//_outbyte(CAN);
//_outbyte(CAN);
//_outbyte(CAN);
//return -2; /* sync error */
//
//start_recv:
//if (trychar == 'C') crc = 1;
//trychar = 0;
//p = xbuff;
//*p++ = c;
//for (i = 0;  i < (bufsz+(crc?1:0)+3); ++i) {
//if ((c = _inbyte(DLY_1S)) < 0) goto reject;
//*p++ = c;
//}
//
//if (xbuff[1] == (unsigned char)(~xbuff[2]) &&
//(xbuff[1] == packetno || xbuff[1] == (unsigned char)packetno-1) &&
//check(crc, &xbuff[3], bufsz)) {
//if (xbuff[1] == packetno)	{
//register int count = destsz - len;
//if (count > bufsz) count = bufsz;
//if (count > 0) {
//memcpy (&dest[len], &xbuff[3], count);
//len += count;
//}
//++packetno;
//retrans = MAXRETRANS+1;
//}
//if (--retrans <= 0) {
//flushinput();
//_outbyte(CAN);
//_outbyte(CAN);
//_outbyte(CAN);
//return -3; /* too many retry error */
//}
//_outbyte(ACK);
//continue;
//}
//reject:
//flushinput();
//_outbyte(NAK);
//}
//}
//

