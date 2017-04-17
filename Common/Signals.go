package Common

var SpecialComStat bool
var SpecialComTast int

const(
	File_Trans =0x10

)

var Ch_ComStreamData chan []byte
var Ch_Ft_start chan int

func Init(){
	Ch_ComStreamData = make(chan []byte,512)
	Ch_Ft_start = make(chan int)
}