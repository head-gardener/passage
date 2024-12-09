package belt

import (
	"crypto/cipher"
)

type beltAEAD struct {
	key  Key
	wrap func(
		out []byte,
		crit []byte,
		open []byte,
		key Key,
		iv IV,
		opt *AEADOpt,
	) (mac MAC, err error)
	unwrap func(
		out []byte,
		crit []byte,
		open []byte,
		mac MAC,
		key Key,
		iv IV,
		opt *AEADOpt,
	) (err error)
}

const MACSize = len(MAC{})

func (b beltAEAD) NonceSize() int { return 64 }

func (b beltAEAD) Overhead() int { return MACSize }

func (b beltAEAD) Open(dst []byte, nonce []byte, ciphertext []byte, additionalData []byte) ([]byte, error) {
	mac := MAC(ciphertext[len(ciphertext)-MACSize:])
	ciphertext = ciphertext[:len(ciphertext)-MACSize]
	res, out := sliceForAppend(dst, len(ciphertext))
	err := b.unwrap(out, ciphertext, additionalData, mac, b.key, IV(nonce), nil)
	return res, err
}

func (b beltAEAD) Seal(dst []byte, nonce []byte, plaintext []byte, additionalData []byte) []byte {
	res, out := sliceForAppend(dst, len(plaintext)+MACSize)
	mac, err := b.wrap(out, plaintext, additionalData, b.key, IV(nonce), nil)
	if err != nil {
		panic(err)
	}
	copy(out[len(plaintext):], mac[:])
	return res
}

func NewCHE(key Key) cipher.AEAD {
	return beltAEAD{key, CHEWrap, CHEUnwrap}
}

func NewDWP(key Key) cipher.AEAD {
	return beltAEAD{key, DWPWrap, DWPUnwrap}
}

func sliceForAppend(in []byte, n int) (head, tail []byte) {
	if total := len(in) + n; cap(in) >= total {
		head = in[:total]
	} else {
		head = make([]byte, total)
		copy(head, in)
	}
	tail = head[len(in):]
	return
}
