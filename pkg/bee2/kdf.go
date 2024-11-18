package bee2

// #cgo LDFLAGS: -lbee2
// #include <stdlib.h>
// #include <bee2/crypto/belt.h>
import "C"

import (
	"fmt"
	"unsafe"
)

type KDFOpt struct {
	passLen int
	saltLen int
	iter    int
}

// Belt key derivation via bee2.
func KDF(
	pass []byte,
	salt []byte,
	opt *KDFOpt,
) (key BeltKey, err error) {
	var passLen int
	if opt != nil && opt.passLen != 0 {
		passLen = opt.passLen
	} else {
		passLen = len(pass)
	}
	var saltLen int
	if opt != nil && opt.saltLen != 0 {
		saltLen = opt.saltLen
	} else {
		saltLen = len(salt)
	}
	var iter int
	if opt != nil && opt.iter != 0 {
		iter = opt.iter
	} else {
		iter = 10000
	}

	if len(pass) == 0 {
		return key, fmt.Errorf("empty pass")
	}
	if len(salt) == 0 {
		return key, fmt.Errorf("empty salt")
	}

	ret := C.beltPBKDF2(
		(*C.octet)(unsafe.Pointer(&key[0])),
		(*C.octet)(unsafe.Pointer(&pass[0])),
		(C.size_t)(passLen),
		(C.size_t)(iter),
		(*C.octet)(unsafe.Pointer(&salt[0])),
		(C.size_t)(saltLen),
	)
	return key, errorMessage(ret)
}
