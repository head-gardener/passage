package crypto

import (
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
	return key
}

func mustDerive(t *testing.T, str []byte) (key BeltKey) {
	key, err := KDF(str, []byte{0}, &KDFOpt{iter: 1})
	if err != nil {
		t.Fatal(err)
	}
	return
}

// Produced by KDF, consumed by everything else. A 32 byte, 256 bit slice.
type BeltKey [32]byte

type CommonOpt struct {
	srcLen int
}
