package vscanapi

/*
#cgo CFLAGS: -I./include
#cgo LDFLAGS: -lvs_can_api
#cgo windows CFLAGS: -DWIN32
#cgo windows amd64 LDFLAGS: -L./lib/win64
#cgo windows 386 LDFLAGS: -L./lib/win32
#cgo linux amd64 LDFLAGS: -L./lib/linux64
#cgo linux 386 LDFLAGS: -L./lib/linux32

#include <stdio.h>
#include <stdlib.h>
#include <stdint.h>
#include <vs_can_api.h>
*/
import "C"
import (
	"fmt"
	"unsafe"
)

// VSCANError is the error returned by VSCAN functions
type VSCANError struct {
	Code int
}

func (e *VSCANError) Error() string {
	buf := make([]byte, 64)

	C.VSCAN_GetErrorString(C.int(e.Code), (*C.char)(unsafe.Pointer(&buf[0])), 64)

	return string(buf)
}

func (e *VSCANError) String() string {
	return e.Error()
}

func (e *VSCANError) Unwrap() error {
	return e
}

type VSCANMode int

const (
	ModeNormal        VSCANMode = C.VSCAN_MODE_NORMAL
	ModeListenOnly    VSCANMode = C.VSCAN_MODE_LISTEN_ONLY
	ModeSelfReception VSCANMode = C.VSCAN_MODE_SELF_RECEPTION
)

type VSCAN struct {
	Handle int
}

func Open(port string, mode VSCANMode) (*VSCAN, error) {
	var handle int

	cport := C.CString(port)
	err := C.VSCAN_Open(cport, C.ulong(mode))

	if err < 0 {
		return nil, &VSCANError{int(err)}
	}

	return &VSCAN{handle}, nil
}

func (v *VSCAN) Close() error {
	err := C.VSCAN_Close(C.int(v.Handle))

	if err < 0 {
		return &VSCANError{int(err)}
	}

	return nil
}

type VSCANMessage struct {
	ID        uint32
	Size      uint8
	Data      [8]uint8
	Flags     uint8
	Timestamp uint16
}

func (v *VSCAN) Read() ([]VSCANMessage, error) {
	msgs := make([]C.VSCAN_MSG, 32)

	var read C.ulong

	err := C.VSCAN_Read(C.int(v.Handle), (*C.VSCAN_MSG)(unsafe.Pointer(&msgs)), C.ulong(len(msgs)), &read)

	if err < 0 {
		return nil, &VSCANError{int(err)}
	}

	var msgsOut []VSCANMessage

	for i := 0; i < int(read); i++ {
		var msg VSCANMessage

		msg.ID = uint32(msgs[i].Id)
		msg.Size = uint8(msgs[i].Size)
		msg.Flags = uint8(msgs[i].Flags)
		msg.Timestamp = uint16(msgs[i].Timestamp)
		for j := 0; j < int(msg.Size); j++ {
			msg.Data[j] = uint8(msgs[i].Data[j])
		}

		msgsOut = append(msgsOut, msg)
	}

	return msgsOut, nil
}

func (v *VSCAN) Write(msg []VSCANMessage) error {
	msgs := make([]C.VSCAN_MSG, len(msg))

	for i := 0; i < len(msg); i++ {
		msgs[i].Id = C.uint(msg[i].ID)
		msgs[i].Size = C.uchar(msg[i].Size)
		msgs[i].Flags = C.uchar(msg[i].Flags)
		msgs[i].Timestamp = C.ushort(msg[i].Timestamp)
		for j := 0; j < int(msg[i].Size); j++ {
			msgs[i].Data[j] = C.uchar(msg[i].Data[j])
		}
	}

	var written C.ulong

	err := C.VSCAN_Write(C.int(v.Handle), (*C.VSCAN_MSG)(unsafe.Pointer(&msgs)), C.ulong(len(msgs)), &written)

	if err < 0 {
		return &VSCANError{int(err)}
	}

	return nil
}

type VSCANSpeed int

const (
	Speed_1M   VSCANSpeed = 8
	Speed_800k VSCANSpeed = 7
	Speed_500k VSCANSpeed = 6
	Speed_250k VSCANSpeed = 5
	Speed_125k VSCANSpeed = 4
	Speed_100k VSCANSpeed = 3
	Speed_50k  VSCANSpeed = 2
	Speed_20k  VSCANSpeed = 1
)

func (v *VSCAN) SetSpeed(speed VSCANSpeed) error {
	err := C.VSCAN_Ioctl(C.int(v.Handle), C.VSCAN_IOCTL_SET_SPEED, unsafe.Pointer(&speed))

	if err < 0 {
		return &VSCANError{int(err)}
	}

	return nil
}

type VSCANApiVersion struct {
	Major    int
	Minor    int
	SubMinor int
}

func (v *VSCANApiVersion) String() string {
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.SubMinor)
}

func GetApiVersion() (VSCANApiVersion, error) {
	var version C.VSCAN_API_VERSION

	err := C.VSCAN_Ioctl(0, C.VSCAN_IOCTL_GET_API_VERSION, unsafe.Pointer(&version))

	if err < 0 {
		return VSCANApiVersion{}, &VSCANError{int(err)}
	}

	return VSCANApiVersion{int(version.Major), int(version.Minor), int(version.SubMinor)}, nil
}
