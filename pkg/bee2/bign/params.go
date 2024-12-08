package bign

// #cgo LDFLAGS: -lbee2
// #include <stdlib.h>
// #include <bee2/crypto/bign.h>
import "C"

import (
	"unsafe"
)

type CurveName string

const (
	L128 CurveName = "1.2.112.0.2.0.34.101.45.3.1"
	L192 CurveName = "1.2.112.0.2.0.34.101.45.3.2"
	L256 CurveName = "1.2.112.0.2.0.34.101.45.3.3"
)

type Params [C.sizeof_bign_params]byte

var (
	P128 Params
	P192 Params
	P256 Params
)

func init() {
	P128, _ = StandardParams(L128)
	P192, _ = StandardParams(L192)
	P256, _ = StandardParams(L256)
}

func (p *Params) toPtr() *C.struct___0 {
	return (*C.struct___0)(unsafe.Pointer(&p[0]))
}

func StandardParams(n CurveName) (p Params, err error) {
	name := C.CString((string)(n))
	defer C.free(unsafe.Pointer(name))
	ret := C.bignParamsStd(
		p.toPtr(),
		name,
	)
	err = errorMessage(ret)
	return
}

func (p *Params) GetL() uint8 {
	return p[0]
}

func (p *Params) Validate() (err error) {
	ret := C.bignParamsVal(
		(*C.struct___0)(unsafe.Pointer(&p[0])),
	)
	err = errorMessage(ret)
	return
}
