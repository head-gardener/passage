package crypto

import (
	"github.com/head-gardener/passage/pkg/bee2"

	"lukechampine.com/uint128"
)

type MAC = bee2.BeltMAC

type Cipher interface {
	Wrap(out []byte, crit []byte, open []byte, mac []byte) error
	Unwrap(out []byte, crit []byte, open []byte, mac bee2.BeltMAC) error
	Inc()
	Finalize()
}

type CHE struct {
	iv  bee2.BeltIV
	key bee2.BeltKey
}

func InitCHE(pass []byte, salt []byte) (che *CHE, err error) {
	key, err := bee2.KDF(pass, salt, nil)
	if err != nil {
		return
	}

	return &CHE{
		key: key,
		iv:  bee2.BeltIV{},
	}, nil
}

func (che *CHE) Wrap(out []byte, crit []byte, open []byte, mac []byte) (err error) {
	m, err := bee2.CHEWrap(out, crit, open, che.key, che.iv, nil)
	if err != nil {
		return
	}

	copy(mac, m[:])
	return nil
}

func (che *CHE) Unwrap(out []byte, crit []byte, open []byte, mac bee2.BeltMAC) (err error) {
	err = bee2.CHEUnwrap(out, crit, open, mac, che.key, che.iv, nil)
	if err != nil {
		return
	}

	return nil
}

func (che *CHE) Inc() {
	uint128.FromBytes(che.iv[:]).AddWrap64(1).PutBytes(che.iv[:])
}

func (che *CHE) Finalize() {
	for i := range len(che.iv) {
		che.iv[i] = 0
	}
	for i := range len(che.key) {
		che.key[i] = 0
	}
}
