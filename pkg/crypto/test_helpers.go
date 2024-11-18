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

func mustBeMAC(t *testing.T, str string) (mac BeltMAC) {
	buf := mustDecode(t, str)
	if len(buf) != 8 {
		t.Fatalf("invalid length %d for an mac", len(buf))
	}
	copy(mac[:], buf)
	return
}

type dataAEAD struct {
	crit []byte
	open []byte
	mac  BeltMAC
}

func compare(a interface{}, b interface{}) bool {
	if reflect.TypeOf(a) != reflect.TypeOf(b) {
		return false
	}

	switch a.(type) {
	case []byte:
		return bytes.Equal(a.([]byte), b.([]byte))
	case *dataAEAD:
		return bytes.Equal(a.(*dataAEAD).crit, b.(*dataAEAD).crit)
	default:
		return reflect.DeepEqual(a, b)
	}
}

type cryptoF[D any, O any] func(out D, src D, key BeltKey, iv BeltIV, opt *O) error
type cryptoI[D any, O any] func(t *testing.T, input D, key BeltKey, iv BeltIV, buf D)
type cryptoC[D any, O any] func(t *testing.T, input D, want D, key BeltKey, iv BeltIV, buf D)

func makeCryptoHelpers[D any, O any](
	opt *O,
	encrypt cryptoF[D, O],
	decrypt cryptoF[D, O],
) (
	identity cryptoI[D, O],
	check cryptoC[D, O],
) {
	return makeCryptoIdentity(opt, encrypt, decrypt), makeCheck(opt, encrypt)
}

func makeCheck[D any, O any](opt *O, encrypt cryptoF[D, O]) cryptoC[D, O] {
	return func(t *testing.T, input D, want D, key BeltKey, iv BeltIV, buf D) {
		fail := func(err any) {
			t.Fatalf(
				"\nerr:   %v\nbuf:   %x\nkey:   %x\ninput: %x\nwant:  %x",
				interface{}(err),
				interface{}(buf),
				interface{}(key),
				interface{}(input),
				interface{}(want),
			)
		}

		if err := encrypt(buf, input, key, iv, opt); err != nil {
			fail(err)
		}
		if !reflect.DeepEqual(buf, want) {
			fail("invalid encryption, buf != want")
		}
	}
}

func makeCryptoIdentity[D any, O any](
	opt *O,
	encrypt cryptoF[D, O],
	decrypt cryptoF[D, O],
) cryptoI[D, O] {
	return func(t *testing.T, input D, key BeltKey, iv BeltIV, buf D) {
		fail := func(err any) {
			t.Fatalf(
				"\nerr:   %v\nbuf:   %x\nkey:   %x\ninput: %x",
				interface{}(err),
				interface{}(buf),
				interface{}(key),
				interface{}(input),
			)
		}

		if err := encrypt(buf, input, key, iv, opt); err != nil {
			fail(err)
		}
		if err := decrypt(buf, buf, key, iv, opt); err != nil {
			fail(err)
		}
		if !compare(buf, input) {
			fail("no decryption, dec != input")
		}
	}
}
