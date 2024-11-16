package crypto

// #cgo LDFLAGS: -lbee2
// #include <stdlib.h>
// #include <bee2/crypto/belt.h>
import "C"

import (
	"fmt"
	"unsafe"
)

// Belt hash via bee2. `out` should be 32 octets long, i.e. 32
func Hash(str []byte, out []byte) (err error) {
	ret := C.beltHash(
		(*C.octet)(unsafe.Pointer(&out[0])),
		unsafe.Pointer(&str[0]),
		(C.size_t)(len(str)),
	)
	if ret != 0 {
		return fmt.Errorf("non-zero return: %v", ret)
	}
	return nil
}
