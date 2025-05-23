package belt

// #cgo LDFLAGS: -lbee2_static
// #include <stdlib.h>
// #include <bee2/crypto/belt.h>
import "C"

import (
	"fmt"
	"unsafe"
)

// Belt electronic codeblock decryption via bee2.
func ECBDecr(
	out []byte,
	src []byte,
	key Key,
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

	ret := C.beltECBDecr(
		unsafe.Pointer(&out[0]),
		unsafe.Pointer(&src[0]),
		(C.size_t)(srcLen),
		(*C.octet)(unsafe.Pointer(&key[0])),
		32,
	)
	return errorMessage(ret)
}

// Belt electronic codeblock encryption via bee2.
func ECBEncr(
	out []byte,
	src []byte,
	key Key,
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

	ret := C.beltECBEncr(
		unsafe.Pointer(&out[0]),
		unsafe.Pointer(&src[0]),
		(C.size_t)(srcLen),
		(*C.octet)(unsafe.Pointer(&key[0])),
		32,
	)
	return errorMessage(ret)
}
