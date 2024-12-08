package belt

// #cgo LDFLAGS: -lbee2
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
