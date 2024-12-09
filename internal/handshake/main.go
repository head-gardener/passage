package handshake

import (
	"crypto/cipher"
	"crypto/rand"
	"encoding/binary"
	"hash"
	"io"

	"github.com/flynn/noise"

	"github.com/head-gardener/passage/pkg/bee2/belt"
	"github.com/head-gardener/passage/pkg/bee2/bign"
)

var BignBeltSuite = noise.NewCipherSuite(Bign128, BeltCHE, BeltHash)

func HandshakeInit(initiator bool, psk []byte) (*noise.HandshakeState, error) {
	return noise.NewHandshakeState(noise.Config{
		CipherSuite:           BignBeltSuite,
		Random:                rand.Reader,
		Pattern:               noise.HandshakeNN,
		Initiator:             initiator,
		PresharedKey:          psk,
		PresharedKeyPlacement: 0,
	})
}

// DF

var Bign128 noise.DHFunc = bignstate{}

type bignstate struct{}

func (bignstate) GenerateKeypair(rng io.Reader) (key noise.DHKey, err error) {
	// NOTE: rng is ignored because passing callbacks to c is too complicated.
	// In practice it's always rand.Reader so it's fine.
	priv, pub, err := bign.GenerateKeypair(&bign.P128)
	if err != nil {
		return
	}
	return noise.DHKey{Private: priv, Public: pub}, nil
}

func (bignstate) DH(privkey, pubkey []byte) ([]byte, error) {
	out := make([]byte, 64)
	err := bign.DiffieHellman(out, privkey, pubkey, &bign.P128)
	return out, err
}

func (bignstate) DHLen() int { return 64 }

func (bignstate) DHName() string { return "Bign128" }

// Cipher

var BeltCHE noise.CipherFunc = beltchefuncstate{}

type beltchefuncstate struct{}

func (b beltchefuncstate) Cipher(k [32]byte) noise.Cipher {
	return beltCHE{belt.NewCHE(k)}
}

func (b beltchefuncstate) CipherName() string { return "BeltCHE" }

type beltCHE struct {
	c cipher.AEAD
}

func (b beltCHE) Decrypt(out []byte, n uint64, ad []byte, ciphertext []byte) ([]byte, error) {
	var nonce [16]byte
	binary.LittleEndian.PutUint64(nonce[8:], n)
	return b.c.Open(out, nonce[:], ciphertext, ad)
}

func (b beltCHE) Encrypt(out []byte, n uint64, ad []byte, plaintext []byte) []byte {
	var nonce [16]byte
	binary.LittleEndian.PutUint64(nonce[8:], n)
	return b.c.Seal(out, nonce[:], plaintext, ad)
}

// Hash

var BeltHash noise.HashFunc = belthashstate{}

type belthashstate struct{}

func (b belthashstate) Hash() hash.Hash {
	return belt.HashInit()
}

func (b belthashstate) HashName() string { return "BeltHash" }
