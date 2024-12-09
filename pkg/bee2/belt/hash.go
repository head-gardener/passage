package belt

// #cgo LDFLAGS: -lbee2
// #include <stdlib.h>
// #include <bee2/crypto/belt.h>
import "C"

import (
	"fmt"
	"hash"
	"unsafe"
)

type hashstate struct {
	state []byte
}

func (h *hashstate) toPtr() unsafe.Pointer {
	return unsafe.Pointer(&h.state[0])
}

func HashInit() hash.Hash {
	h := &hashstate{make([]byte, C.beltHash_keep())}
	C.beltHashStart(h.toPtr())
	return h
}

func (h *hashstate) BlockSize() int { return 32 }

func (h *hashstate) Size() int { return 32 }

func (h *hashstate) Reset() {
	C.beltHashStart(h.toPtr())
}

func (h *hashstate) Sum(b []byte) []byte {
	buf := make([]byte, h.Size())
	C.beltHashStepG(
		(*C.octet)(unsafe.Pointer(&buf[0])),
		h.toPtr(),
	)
	return append(b, buf[:h.Size()]...)
}

func (h *hashstate) Write(p []byte) (int, error) {
	if len(p) == 0 {
		// NOTE: is this expected?
		p = make([]byte, h.BlockSize())
	}
	C.beltHashStepH(
		unsafe.Pointer(&p[0]),
		C.size_t(len(p)),
		h.toPtr(),
	)
	return len(p), nil
}

// Belt hash via bee2. `out` should be 32 bytes long
func Hash(
	out []byte,
	src []byte,
	opt *CommonOpt,
) (err error) {
	var srcLen int
	if opt != nil && opt.srcLen != 0 {
		srcLen = opt.srcLen
	} else {
		srcLen = len(src)
	}

	if len(out) == 0 {
		return fmt.Errorf("empty out")
	}
	if len(src) == 0 {
		return fmt.Errorf("empty src")
	}

	ret := C.beltHash(
		(*C.octet)(unsafe.Pointer(&out[0])),
		unsafe.Pointer(&src[0]),
		(C.size_t)(srcLen),
	)
	return errorMessage(ret)
}
