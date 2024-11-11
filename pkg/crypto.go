package pkg

// #cgo LDFLAGS: -lbee2
// #include <stdlib.h>
// #include <bee2/crypto/belt.h>
import "C"

import (
	"fmt"
	"unsafe"
)

func Encrypt(str []byte) ([]byte, error) {
	buf := C.malloc(4096)
	defer C.free(unsafe.Pointer(buf))
	out := make([]byte, C.sizeof_octet*32)

	fmt.Printf("%x\n", out)

	C.beltHashStart(buf)
	C.beltHashStepH(unsafe.Pointer(&str[0]), 0, buf)
	C.beltHashStepG((*C.octet)(unsafe.Pointer(&out[0])), buf)

	fmt.Printf("%x\n", out)

	return out, nil
}
