package crypto

import (
	"bytes"
	"encoding/hex"
	"math/rand"
	"reflect"
	"testing"
	"testing/quick"
)

const maxSize = 100000

func conf(maxSize int, argn int) quick.Config {
	buf := make([]byte, maxSize*argn)
	return quick.Config{
		Values: func(args []reflect.Value, rand *rand.Rand) {
			for i := range args {
				offset := maxSize * i
				n := rand.Intn(maxSize) + 1
				slice := buf[offset : offset+n]
				rand.Read(slice)
				args[i] = reflect.ValueOf(slice)
			}
		},
	}
}

func mustDecode(t *testing.T, str string) (res []byte) {
	res, err := hex.DecodeString(str)
	if err != nil {
		t.Fatal(err)
	}
	return
}

func mustBeKey(t *testing.T, str string) (key BeltKey) {
	keyBuf := mustDecode(t, str)
	if len(keyBuf) != 32 {
		t.Fatalf("invalid length %d for belt key", len(keyBuf))
	}
	copy(key[:], keyBuf)
	return
}

func mustDerive(t *testing.T, str []byte) (key BeltKey) {
	key, err := KDF(str, []byte{0}, &KDFOpt{iter: 1})
	if err != nil {
		t.Fatal(err)
	}
	return
}

func mustBeIV(t *testing.T, str string) (iv BeltIV) {
	ivBuf := mustDecode(t, str)
	if len(ivBuf) != 16 {
		t.Fatalf("invalid length %d for an iv", len(ivBuf))
	}
	copy(iv[:], ivBuf)
	return
}

func mustContainIV(t *testing.T, str []byte) (iv BeltIV) {
	if len(str) < 16 {
		t.Fatalf("insufficient length %d to form an iv", len(str))
	}
	copy(iv[:], str[:16])
	return
}

func makeCryptoIdentity(
	opt *CommonOpt,
	encrypt func(out []byte, src []byte, key BeltKey, iv BeltIV, opt *CommonOpt) error,
	decrypt func(out []byte, src []byte, key BeltKey, iv BeltIV, opt *CommonOpt) error,
) func(t *testing.T, input []byte, key BeltKey, iv BeltIV, enc []byte, dec []byte) {
	return func(t *testing.T, input []byte, key BeltKey, iv BeltIV, enc []byte, dec []byte) {
		if err := encrypt(enc, input, key, iv, opt); err != nil {
			t.Fatal(err)
		}
		if err := decrypt(dec, enc, key, iv, opt); err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(dec, input) {
			t.Fatalf("no decryption:\n%x + %x ->\n%x, not\n%x", enc, key, dec, input)
		}

		// in-place enc/dec
		if err := encrypt(dec, dec, key, iv, opt); err != nil {
			t.Fatal(err)
		}
		if err := decrypt(dec, dec, key, iv, opt); err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(dec, input) {
			t.Fatalf("in-place no decryption:\n%x + %x ->\n%x, not\n%x", enc, key, dec, input)
		}
	}
}

// Produced by KDF, consumed by everything else. Same key shouldn't be used in different
// algorithms. A 32 byte, 256 bit slice.
type BeltKey [32]byte

// Must be unique for every session using a single key. A 16 byte, 128 bit slice.
type BeltIV [16]byte

type CommonOpt struct {
	srcLen int
}
