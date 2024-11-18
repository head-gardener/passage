package crypto

// #cgo LDFLAGS: -lbee2
// #include <stdlib.h>
// #include <bee2/crypto/belt.h>
import "C"

import (
	"unsafe"
)

// Belt hash via bee2. `out` should be 32 bytes long
func Hash(
	out []byte,
	src []byte,
	opt *CommonOpt,
) (err error) {
	var srcLen int
	if opt != nil && opt.srcLen != 0 {
		srcLen = opt.srcLen
	} else {
		srcLen = len(src)
	}

	ret := C.beltHash(
		(*C.octet)(unsafe.Pointer(&out[0])),
		unsafe.Pointer(&src[0]),
		(C.size_t)(srcLen),
	)
	return errorMessage(ret)
}
