package crypto

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

// Belt key derivation via bee2. `out` should be 32 bytes long
func KDF(
	out []byte,
	pass []byte,
	salt []byte,
	opt *KDFOpt,
) (err error) {
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

	ret := C.beltPBKDF2(
		(*C.octet)(unsafe.Pointer(&out[0])),
		(*C.octet)(unsafe.Pointer(&pass[0])),
		(C.size_t)(passLen),
		(C.size_t)(iter),
		(*C.octet)(unsafe.Pointer(&salt[0])),
		(C.size_t)(saltLen),
	)
	if ret != 0 {
		return fmt.Errorf("non-zero return: %v", ret)
	}
	return nil
}

type HMACOpt struct {
	srcLen int
	keyLen int
}

// Belt HMAC via bee2. `out` should be 32 bytes long
func HMAC(
	out []byte,
	src []byte,
	key []byte,
	opt *HMACOpt,
) (err error) {
	var srcLen int
	if opt != nil && opt.srcLen != 0 {
		srcLen = opt.srcLen
	} else {
		srcLen = len(src)
	}
	var keyLen int
	if opt != nil && opt.srcLen != 0 {
		keyLen = opt.srcLen
	} else {
		keyLen = len(key)
	}

	ret := C.beltHMAC(
		(*C.octet)(unsafe.Pointer(&out[0])),
		unsafe.Pointer(&src[0]),
		(C.size_t)(srcLen),
		(*C.octet)(unsafe.Pointer(&key[0])),
		(C.size_t)(keyLen),
	)
	if ret != 0 {
		return fmt.Errorf("non-zero return: %v", ret)
	}
	return nil
}

type HashOpt struct {
	srcLen int
}

// Belt hash via bee2. `out` should be 32 bytes long
func Hash(
	out []byte,
	src []byte,
	opt *HashOpt,
) (err error) {
	var srcLen int
	if opt != nil && opt.srcLen != 0 {
		srcLen = opt.srcLen
	} else {
		srcLen = len(src)
	}

	ret := C.beltHash(
		(*C.octet)(unsafe.Pointer(&out[0])),
		unsafe.Pointer(&src[0]),
		(C.size_t)(srcLen),
	)
	if ret != 0 {
		return fmt.Errorf("non-zero return: %v", ret)
	}
	return nil
}
