package bign

import (
	"testing"
	"testing/quick"

	"github.com/head-gardener/passage/pkg/bee2/belt"
)

func TestSignature(t *testing.T) {
	sig := make([]byte, 48)
	hash := make([]byte, 32)
	priv, pub, err := GenerateKeypair(&P128)
	if err != nil {
		t.Fatal(err)
	}

	f := func(b []byte) bool {
		if len(b) == 0 {
			return true
		}

		err := belt.Hash(hash, b, nil)
		if err != nil {
			t.Error(err)
			return false
		}
		err = Sign(sig, hash, priv, &P128, nil)
		if err != nil {
			t.Error(err)
			return false
		}
		err = Verify(sig, hash, pub, &P128, nil)
		if err != nil {
			t.Error(err)
			return false
		}
		return true
	}

	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}
