package bee2

// #cgo LDFLAGS: -lbee2
// #include <stdlib.h>
// #include <bee2/crypto/belt.h>
import "C"

import (
	"fmt"
	"unsafe"
)

// Belt AEAD (DWP mode) unwrapping via bee2.
func DWPUnwrap(
	out []byte,
	crit []byte,
	open []byte,
	mac BeltMAC,
	key BeltKey,
	iv BeltIV,
	opt *AEADOpt,
) (err error) {
	var critLen int
	if opt != nil && opt.critLen != 0 {
		critLen = opt.critLen
	} else {
		critLen = len(crit)
	}
	var openLen int
	if open == nil {
		open = crit
		openLen = 0
	} else if opt != nil && opt.openLen != 0 {
		openLen = opt.openLen
	} else {
		openLen = len(open)
	}

	if len(out) == 0 {
		return fmt.Errorf("empty out")
	}
	if len(crit) == 0 {
		return fmt.Errorf("empty crit")
	}
	if len(open) == 0 {
		return fmt.Errorf("empty open")
	}

	ret := C.beltDWPUnwrap(
		unsafe.Pointer(&out[0]),
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

// Belt AEAD (DWP mode) wrapping via bee2.
func DWPWrap(
	out []byte,
	crit []byte,
	open []byte,
	key BeltKey,
	iv BeltIV,
	opt *AEADOpt,
) (mac BeltMAC, err error) {
	var critLen int
	if opt != nil && opt.critLen != 0 {
		critLen = opt.critLen
	} else {
		critLen = len(crit)
	}
	var openLen int
	if open == nil {
		open = crit
		openLen = 0
	} else if opt != nil && opt.openLen != 0 {
		openLen = opt.openLen
	} else {
		openLen = len(open)
	}

	if len(out) == 0 {
		return mac, fmt.Errorf("empty out")
	}
	if len(crit) == 0 {
		return mac, fmt.Errorf("empty crit")
	}
	if len(open) == 0 {
		return mac, fmt.Errorf("empty open")
	}

	ret := C.beltDWPWrap(
		unsafe.Pointer(&out[0]),
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
