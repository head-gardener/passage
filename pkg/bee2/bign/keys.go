package bign

// #cgo LDFLAGS: -lbee2_static
// #include <stdlib.h>
// #include <bee2/crypto/bign.h>
// #include "urand_gen.h"
import "C"

import (
	"unsafe"
)

func GenerateKeypair(
	p *Params,
) (priv []byte, pub []byte, err error) {
	l := p.GetL()
	priv = make([]byte, l/4)
	pub = make([]byte, l/2)
	ret := C.bignKeypairGen(
		(*C.octet)(unsafe.Pointer(&priv[0])),
		(*C.octet)(unsafe.Pointer(&pub[0])),
		p.toPtr(),
		(*[0]byte)(C.urand_gen),
		nil,
	)
	err = errorMessage(ret)
	return
}

func ValidateKeypair(
	pub []byte,
	priv []byte,
	p *Params,
) (err error) {
	ret := C.bignKeypairVal(
		p.toPtr(),
		(*C.octet)(unsafe.Pointer(&priv[0])),
		(*C.octet)(unsafe.Pointer(&pub[0])),
	)
	return errorMessage(ret)
}
