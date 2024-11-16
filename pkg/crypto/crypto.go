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
func Hash(out []byte, str []byte) (err error) {
	return HashWithLength(out, str, len(str))
}

// Same as Hash but input length is explicit. Useful for in-place hashing
func HashWithLength(out []byte, str []byte, length int) (err error) {
	ret := C.beltHash(
		(*C.octet)(unsafe.Pointer(&out[0])),
		unsafe.Pointer(&str[0]),
		(C.size_t)(length),
	)
	if ret != 0 {
		return fmt.Errorf("non-zero return: %v", ret)
	}
	return nil
}
