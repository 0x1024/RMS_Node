package Serial_Srv

import (
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

var _ unsafe.Pointer

var (
	modadvapi32 = windows.NewLazySystemDLL("advapi32.dll")
	modkernel32 = windows.NewLazySystemDLL("kernel32.dll")

	procRegEnumValueW       = modadvapi32.NewProc("RegEnumValueW")
	procGetCommState        = modkernel32.NewProc("GetCommState")
	procSetCommState        = modkernel32.NewProc("SetCommState")
	procSetCommTimeouts     = modkernel32.NewProc("SetCommTimeouts")
	procEscapeCommFunction  = modkernel32.NewProc("EscapeCommFunction")
	procGetCommModemStatus  = modkernel32.NewProc("GetCommModemStatus")
	procCreateEventW        = modkernel32.NewProc("CreateEventW")
	procResetEvent          = modkernel32.NewProc("ResetEvent")
	procGetOverlappedResult = modkernel32.NewProc("GetOverlappedResult")
)

// PortError is a platform independent error type for serial ports
type PortError struct {
	code     PortErrorCode
	causedBy error
}

// PortErrorCode is a code to easily identify the type of error
type PortErrorCode int

const (
	// PortBusy the serial port is already in used by another process
	PortBusy PortErrorCode = iota
	// PortNotFound the requested port doesn't exist
	PortNotFound
	// InvalidSerialPort the requested port is not a serial port
	InvalidSerialPort
	// PermissionDenied the user doesn't have enough priviledges
	PermissionDenied
	// InvalidSpeed the requested speed is not valid or not supported
	InvalidSpeed
	// InvalidDataBits the number of data bits is not valid or not supported
	InvalidDataBits
	// InvalidParity the selected parity is not valid or not supported
	InvalidParity
	// InvalidStopBits the selected number of stop bits is not valid or not supported
	InvalidStopBits
	// ErrorEnumeratingPorts an error occurred while listing serial port
	ErrorEnumeratingPorts
	// PortClosed the port has been closed while the operation is in progress
	PortClosed
	// FunctionNotImplemented the requested function is not implemented
	FunctionNotImplemented
)

// GetPortsList retrieve the list of available serial ports
func GetPortsList() ([]string, error) {
	return nativeGetPortsList()
}

func nativeGetPortsList() ([]string, error) {
	subKey, err := syscall.UTF16PtrFromString("HARDWARE\\DEVICEMAP\\SERIALCOMM\\")
	if err != nil {
		return nil, &PortError{code: ErrorEnumeratingPorts}
	}

	var h syscall.Handle
	if syscall.RegOpenKeyEx(syscall.HKEY_LOCAL_MACHINE, subKey, 0, syscall.KEY_READ, &h) != nil {
		return nil, &PortError{code: ErrorEnumeratingPorts}
	}
	defer syscall.RegCloseKey(h)

	var valuesCount uint32
	if syscall.RegQueryInfoKey(h, nil, nil, nil, nil, nil, nil, &valuesCount, nil, nil, nil, nil) != nil {
		return nil, &PortError{code: ErrorEnumeratingPorts}
	}

	list := make([]string, valuesCount)
	for i := range list {
		var data [1024]uint16
		dataSize := uint32(len(data))
		var name [1024]uint16
		nameSize := uint32(len(name))
		if regEnumValue(h, uint32(i), &name[0], &nameSize, nil, nil, &data[0], &dataSize) != nil {
			return nil, &PortError{code: ErrorEnumeratingPorts}
		}
		list[i] = syscall.UTF16ToString(data[:])
	}
	return list, nil
}

func regEnumValue(key syscall.Handle, index uint32, name *uint16, nameLen *uint32, reserved *uint32, class *uint16, value *uint16, valueLen *uint32) (regerrno error) {
	r0, _, _ := syscall.Syscall9(procRegEnumValueW.Addr(), 8, uintptr(key), uintptr(index), uintptr(unsafe.Pointer(name)), uintptr(unsafe.Pointer(nameLen)), uintptr(unsafe.Pointer(reserved)), uintptr(unsafe.Pointer(class)), uintptr(unsafe.Pointer(value)), uintptr(unsafe.Pointer(valueLen)), 0)
	if r0 != 0 {
		regerrno = syscall.Errno(r0)
	}
	return
}

// Error returns the complete error code with details on the cause of the error
func (e PortError) Error() string {
	if e.causedBy != nil {
		return e.EncodedErrorString() + ": " + e.causedBy.Error()
	}
	return e.EncodedErrorString()
}

// Code returns an identifier for the kind of error occurred
func (e PortError) Code() PortErrorCode {
	return e.code
}

// EncodedErrorString returns a string explaining the error code
func (e PortError) EncodedErrorString() string {
	switch e.code {
	case PortBusy:
		return "Serial port busy"
	case PortNotFound:
		return "Serial port not found"
	case InvalidSerialPort:
		return "Invalid serial port"
	case PermissionDenied:
		return "Permission denied"
	case InvalidSpeed:
		return "Port speed invalid or not supported"
	case InvalidDataBits:
		return "Port data bits invalid or not supported"
	case InvalidParity:
		return "Port parity invalid or not supported"
	case InvalidStopBits:
		return "Port stop bits invalid or not supported"
	case ErrorEnumeratingPorts:
		return "Could not enumerate serial ports"
	case PortClosed:
		return "Port has been closed"
	case FunctionNotImplemented:
		return "Function not implemented"
	default:
		return "Other error"
	}
}
