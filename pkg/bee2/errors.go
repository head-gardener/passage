package bee2

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
		return fmt.Errorf("Unknown error %d", uintptr(unsafe.Pointer(ptr)))
	}

	msg := C.GoString(ptr)
	return errors.New(msg)
}