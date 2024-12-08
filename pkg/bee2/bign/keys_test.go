package bign

import (
	"testing"
	"testing/quick"
)

func TestKeyGeneration(t *testing.T) {
	f := func() bool {
		priv, pub, err := GenerateKeypair(&P128)
		if err != nil {
			t.Fatal(err)
		}
		err = ValidateKeypair(pub, priv, &P128)
		if err != nil {
			t.Errorf("priv: %x\npub: %x\nerr: %v\n", priv, pub, err)
			return false
		}
		return true
	}

	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}
