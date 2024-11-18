package crypto

// #cgo LDFLAGS: -lbee2
// #include <stdlib.h>
// #include <bee2/crypto/belt.h>
import "C"

import (
	"unsafe"
)

type HMACOpt struct {
	srcLen int
	keyLen int
}

// Belt HMAC via bee2. `out` should be 32 bytes long
func HMAC(
	out []byte,
	src []byte,
	key []byte,
	opt *HMACOpt,
) (err error) {
	var srcLen int
	if opt != nil && opt.srcLen != 0 {
		srcLen = opt.srcLen
	} else {
		srcLen = len(src)
	}
	var keyLen int
	if opt != nil && opt.srcLen != 0 {
		keyLen = opt.srcLen
	} else {
		keyLen = len(key)
	}

	ret := C.beltHMAC(
		(*C.octet)(unsafe.Pointer(&out[0])),
		unsafe.Pointer(&src[0]),
		(C.size_t)(srcLen),
		(*C.octet)(unsafe.Pointer(&key[0])),
		(C.size_t)(keyLen),
	)
	return errorMessage(ret)
}
