package bign

// #cgo LDFLAGS: -lbee2_static
// #include <stdlib.h>
// #include <bee2/crypto/bign.h>
import "C"

import (
	"fmt"
	"unsafe"
)

func DiffieHellman(
	out []byte,
	priv []byte,
	pub []byte,
	p *Params,
) (err error) {
	if len(out) == 0 {
		return fmt.Errorf("empty out")
	}
	if len(out) > int(p.GetL()/2) {
		return fmt.Errorf(
			"output key length %v exceeds maximum size allowed by security level %v: %v",
			len(out),
			p.GetL(),
			p.GetL()/2,
		)
	}

	ret := C.bignDH(
		(*C.octet)(unsafe.Pointer(&out[0])),
		p.toPtr(),
		(*C.octet)(unsafe.Pointer(&priv[0])),
		(*C.octet)(unsafe.Pointer(&pub[0])),
		(C.size_t)(len(out)),
	)
	return errorMessage(ret)
}
