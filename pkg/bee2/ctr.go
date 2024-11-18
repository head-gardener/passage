package bee2

// #cgo LDFLAGS: -lbee2
// #include <stdlib.h>
// #include <bee2/crypto/belt.h>
import "C"

import (
	"unsafe"
)

// Belt counter encryption via bee2.
func CTR(
	out []byte,
	src []byte,
	key BeltKey,
	iv BeltIV,
	opt *CommonOpt,
) (err error) {
	var srcLen int
	if opt != nil && opt.srcLen != 0 {
		srcLen = opt.srcLen
	} else {
		srcLen = len(src)
	}

	ret := C.beltCTR(
		unsafe.Pointer(&out[0]),
		unsafe.Pointer(&src[0]),
		(C.size_t)(srcLen),
		(*C.octet)(unsafe.Pointer(&key[0])),
		32,
		(*C.octet)(unsafe.Pointer(&iv[0])),
	)
	return errorMessage(ret)
}
