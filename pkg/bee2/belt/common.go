package belt

// #cgo LDFLAGS: -lbee2_static
// #include <stdlib.h>
// #include <bee2/core/err.h>
import "C"

import (
	"errors"
	"fmt"
	"unsafe"
)

func errorMessage(code C.err_t) error {
	if code == 0 {
		return nil
	}

	ptr := C.errMsg((C.err_t)(code))
	if uintptr(unsafe.Pointer(ptr)) == 0 {
		return fmt.Errorf("unknown error %d", uintptr(unsafe.Pointer(ptr)))
	}

	msg := C.GoString(ptr)
	return errors.New(msg)
}

// Produced by KDF, consumed by everything else. Same key shouldn't be used in different
// algorithms. A 32 byte, 256 bit slice.
type Key [32]byte

// Must be unique for every session using a single key. A 16 byte, 128 bit slice.
type IV [16]byte

// Produced and consumed by AEAD functions. A 8 byte, 64 bit slice.
type MAC [8]byte

type CommonOpt struct {
	srcLen int
}

type AEADOpt struct {
	critLen int
	openLen int
}

func prepareOptsAEAD(
	out []byte,
	iv IV,
	c []byte,
	o []byte,
	opt *AEADOpt,
) (outPtr unsafe.Pointer, crit []byte, critLen int, open []byte, openLen int, err error) {
	if opt != nil && opt.critLen != 0 {
		critLen = opt.critLen
	} else {
		critLen = len(c)
	}

	if len(c) == 0 {
		crit = iv[:]
	} else {
		crit = c
	}

	if opt != nil && opt.openLen != 0 {
		openLen = opt.openLen
	} else {
		openLen = len(o)
	}

	if len(o) == 0 {
		open = iv[:]
	} else {
		open = o
	}

	if len(out) == 0 {
		if critLen != 0 {
			return nil, nil, 0, nil, 0, fmt.Errorf("empty out with unempty crit")
		} else {
			outPtr = unsafe.Pointer(&out)
		}
	} else {
		outPtr = unsafe.Pointer(&out[0])
	}

	return
}
