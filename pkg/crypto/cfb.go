package crypto

// #cgo LDFLAGS: -lbee2
// #include <stdlib.h>
// #include <bee2/crypto/belt.h>
import "C"

import (
	"fmt"
	"unsafe"
)

// Belt electronic codeblock decryption via bee2.
func CFBDecr(
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

	ret := C.beltCFBDecr(
		unsafe.Pointer(&out[0]),
		unsafe.Pointer(&src[0]),
		(C.size_t)(srcLen),
		(*C.octet)(unsafe.Pointer(&key[0])),
		32,
		(*C.octet)(unsafe.Pointer(&iv[0])),
	)
	if ret != 0 {
		return fmt.Errorf("non-zero return: %v", ret)
	}
	return nil
}

// Belt electronic codeblock encryption via bee2.
func CFBEncr(
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

	ret := C.beltCFBEncr(
		unsafe.Pointer(&out[0]),
		unsafe.Pointer(&src[0]),
		(C.size_t)(srcLen),
		(*C.octet)(unsafe.Pointer(&key[0])),
		32,
		(*C.octet)(unsafe.Pointer(&iv[0])),
	)
	if ret != 0 {
		return fmt.Errorf("non-zero return: %v", ret)
	}
	return nil
}
