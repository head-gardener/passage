package belt

import (
	"bytes"
	"crypto/cipher"
	"crypto/rand"
	"testing"
)

func TestAEAD(t *testing.T) {
	key := Key{}
	rand.Read(key[:])
	ciphers := []cipher.AEAD{NewCHE(key), NewDWP(key)}
	nonce := IV{}
	rand.Read(nonce[:])
	plaintext := []byte("plaintext")
	additionalData := []byte("additionalData")
	// buf := make([]byte, len(plaintext) + ciphers[0].Overhead())
	for i := range ciphers {
		enc := ciphers[i].Seal(nil, nonce[:], plaintext, additionalData)
		res, err := ciphers[i].Open(nil, nonce[:], enc, additionalData)
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(res, plaintext) {
			t.Fatalf("unseal error: %x vs %x", res, plaintext)
		}

		buf := []byte("wwwwwwww")
		enc = ciphers[i].Seal(buf, nonce[:], plaintext, additionalData)
		if !bytes.Equal(enc[:len(buf)], buf) || len(enc) != len(buf)+ciphers[i].Overhead()+len(plaintext) {
			t.Fatalf(
				"malformed seal: %x vs %x + enc(%x) + [%v]",
				enc,
				buf,
				plaintext,
				ciphers[i].Overhead(),
			)
		}
		res, err = ciphers[i].Open(nil, nonce[:], enc[len(buf):], additionalData)
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(res, plaintext) {
			t.Fatalf("complex unseal error: %x vs %x", res, plaintext)
		}
	}
}
