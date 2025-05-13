package belt

// #cgo LDFLAGS: -lbee2_static
// #include <stdlib.h>
// #include <bee2/crypto/belt.h>
import "C"

import (
	"unsafe"
)

// Belt AEAD (Counter-Hash-Encrypt mode) unwrapping via bee2.
func CHEUnwrap(
	out []byte,
	crit []byte,
	open []byte,
	mac MAC,
	key Key,
	iv IV,
	opt *AEADOpt,
) (err error) {
	outPtr, crit, critLen, open, openLen, err := prepareOptsAEAD(out, iv, crit, open, opt)
	if err != nil {
		return err
	}

	ret := C.beltCHEUnwrap(
		outPtr,
		unsafe.Pointer(&crit[0]),
		(C.size_t)(critLen),
		unsafe.Pointer(&open[0]),
		(C.size_t)(openLen),
		(*C.octet)(unsafe.Pointer(&mac[0])),
		(*C.octet)(unsafe.Pointer(&key[0])),
		32,
		(*C.octet)(unsafe.Pointer(&iv[0])),
	)
	return errorMessage(ret)
}

// Belt AEAD (Counter-Hash-Encrypt mode) wrapping via bee2.
func CHEWrap(
	out []byte,
	crit []byte,
	open []byte,
	key Key,
	iv IV,
	opt *AEADOpt,
) (mac MAC, err error) {
	outPtr, crit, critLen, open, openLen, err := prepareOptsAEAD(out, iv, crit, open, opt)
	if err != nil {
		return mac, err
	}

	ret := C.beltCHEWrap(
		outPtr,
		(*C.octet)(unsafe.Pointer(&mac[0])),
		unsafe.Pointer(&crit[0]),
		(C.size_t)(critLen),
		unsafe.Pointer(&open[0]),
		(C.size_t)(openLen),
		(*C.octet)(unsafe.Pointer(&key[0])),
		32,
		(*C.octet)(unsafe.Pointer(&iv[0])),
	)
	return mac, errorMessage(ret)
}
