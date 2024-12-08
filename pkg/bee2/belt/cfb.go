package belt

// #cgo LDFLAGS: -lbee2
// #include <stdlib.h>
// #include <bee2/crypto/belt.h>
import "C"

import (
	"fmt"
	"unsafe"
)

// Belt cipher feedback decryption via bee2.
func CFBDecr(
	out []byte,
	src []byte,
	key Key,
	iv IV,
	opt *CommonOpt,
) (err error) {
	var srcLen int
	if opt != nil && opt.srcLen != 0 {
		srcLen = opt.srcLen
	} else {
		srcLen = len(src)
	}

	if len(out) == 0 {
		return fmt.Errorf("empty out")
	}
	if len(src) == 0 {
		return fmt.Errorf("empty src")
	}

	ret := C.beltCFBDecr(
		unsafe.Pointer(&out[0]),
		unsafe.Pointer(&src[0]),
		(C.size_t)(srcLen),
		(*C.octet)(unsafe.Pointer(&key[0])),
		32,
		(*C.octet)(unsafe.Pointer(&iv[0])),
	)
	return errorMessage(ret)
}

// Belt cipher feedback encryption via bee2.
func CFBEncr(
	out []byte,
	src []byte,
	key Key,
	iv IV,
	opt *CommonOpt,
) (err error) {
	var srcLen int
	if opt != nil && opt.srcLen != 0 {
		srcLen = opt.srcLen
	} else {
		srcLen = len(src)
	}

	if len(out) == 0 {
		return fmt.Errorf("empty out")
	}
	if len(src) == 0 {
		return fmt.Errorf("empty src")
	}

	ret := C.beltCFBEncr(
		unsafe.Pointer(&out[0]),
		unsafe.Pointer(&src[0]),
		(C.size_t)(srcLen),
		(*C.octet)(unsafe.Pointer(&key[0])),
		32,
		(*C.octet)(unsafe.Pointer(&iv[0])),
	)
	return errorMessage(ret)
}
