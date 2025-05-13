package bign

// #cgo LDFLAGS: -lbee2_static
// #include <stdlib.h>
// #include <bee2/crypto/bign.h>
import "C"

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"unsafe"
)

type SigOpt struct {
	oid_der []byte
	t       []byte
}

var oid_belt_hash []byte

func init() {
	oid_belt_hash, _ = hex.DecodeString("06092a7000020022651f51")
}

func Verify(
	sig []byte,
	hash []byte,
	key []byte,
	p *Params,
	opts *SigOpt,
) (err error) {
	if len(sig) == 0 {
		return fmt.Errorf("empty sig")
	}
	if len(hash) == 0 {
		return fmt.Errorf("empty hash")
	}
	if len(key) == 0 {
		return fmt.Errorf("empty key")
	}

	var (
		oid_der []byte
		oid_len int
	)

	if opts != nil && opts.t != nil {
		oid_der = opts.oid_der
	} else {
		// fine if neither side expects STB compliancy
		oid_der = oid_belt_hash
	}
	oid_len = len(oid_der)

	ret := C.bignVerify(
		p.toPtr(),
		(*C.octet)(unsafe.Pointer(&oid_der[0])),
		(C.size_t)(oid_len),
		(*C.octet)(unsafe.Pointer(&hash[0])),
		(*C.octet)(unsafe.Pointer(&sig[0])),
		(*C.octet)(unsafe.Pointer(&key[0])),
	)
	return errorMessage(ret)
}

func Sign(
	out []byte,
	hash []byte,
	key []byte,
	p *Params,
	opts *SigOpt,
) (err error) {
	if len(out) == 0 {
		return fmt.Errorf("empty out")
	}
	if len(hash) == 0 {
		return fmt.Errorf("empty hash")
	}
	if len(key) == 0 {
		return fmt.Errorf("empty key")
	}

	var t []byte
	if opts != nil && opts.t != nil {
		t = opts.t
	} else {
		t = make([]byte, 32)
		_, err = rand.Read(t)
		if err != nil {
			return
		}
	}

	var (
		oid_der []byte
		oid_len int
	)

	if opts != nil && opts.t != nil {
		oid_der = opts.oid_der
	} else {
		// fine if neither side expects STB compliancy
		oid_der = oid_belt_hash
	}
	oid_len = len(oid_der)

	ret := C.bignSign2(
		(*C.octet)(unsafe.Pointer(&out[0])),
		p.toPtr(),
		(*C.octet)(unsafe.Pointer(&oid_der[0])),
		(C.size_t)(oid_len),
		(*C.octet)(unsafe.Pointer(&hash[0])),
		(*C.octet)(unsafe.Pointer(&key[0])),
		unsafe.Pointer(&t[0]),
		(C.size_t)(len(t)),
	)
	return errorMessage(ret)
}
