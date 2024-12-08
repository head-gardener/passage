package crypto

import (
	"bytes"
	"testing"
	"testing/quick"

	"github.com/head-gardener/passage/pkg/bee2/belt"
)

func TestCiphers(t *testing.T) {
	pass := []byte("pass")
	salt := []byte("salt")
	cipher, err := InitCHE(pass, salt)
	if err != nil {
		t.Fatalf("initialization error: %v", err)
	}

	f := func(b []byte) (ok bool) {
		if len(b) == 0 {
			return true
		}

		ok = false
		buf := make([]byte, len(b))
		mac := belt.MAC{}

		err = cipher.Wrap(buf, b, nil, mac[:])
		if err != nil {
			t.Fatalf("error wrapping: %v", err)
		}
		err = cipher.Unwrap(buf, buf, nil, mac)
		if err != nil {
			t.Fatalf("error unwrapping: %v", err)
		}

		if !bytes.Equal(b, buf) {
			return
		}

		ok = true
		return
	}

	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}
