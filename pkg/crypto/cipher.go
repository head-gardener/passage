package crypto

import (
	"github.com/head-gardener/passage/pkg/bee2"

	"lukechampine.com/uint128"
)

type MAC = bee2.BeltMAC

type Cipher interface {
	Wrap(out []byte, crit []byte, open []byte, mac []byte) error
	Unwrap(out []byte, crit []byte, open []byte, mac bee2.BeltMAC) error
	Finalize()
}

type CHE struct {
	ivu bee2.BeltIV
	ivw bee2.BeltIV
	key bee2.BeltKey
}

func inc128(i *bee2.BeltIV) {
	uint128.FromBytes(i[:]).AddWrap64(1).PutBytes(i[:])
}

func InitCHE(pass []byte, salt []byte) (che *CHE, err error) {
	key, err := bee2.KDF(pass, salt, nil)
	if err != nil {
		return
	}

	return &CHE{
		key: key,
		ivu: bee2.BeltIV{},
		ivw: bee2.BeltIV{},
	}, nil
}

func (che *CHE) Wrap(out []byte, crit []byte, open []byte, mac []byte) (err error) {
	m, err := bee2.CHEWrap(out, crit, open, che.key, che.ivw, nil)
	if err != nil {
		return
	}

	copy(mac, m[:])
	inc128(&che.ivw)
	return nil
}

func (che *CHE) Unwrap(out []byte, crit []byte, open []byte, mac bee2.BeltMAC) (err error) {
	err = bee2.CHEUnwrap(out, crit, open, mac, che.key, che.ivu, nil)
	if err != nil {
		return
	}

	inc128(&che.ivu)

	return nil
}

func (che *CHE) Finalize() {
	for i := range len(che.ivw) {
		che.ivw[i] = 0
	}
	for i := range len(che.ivu) {
		che.ivu[i] = 0
	}
	for i := range len(che.key) {
		che.key[i] = 0
	}
}
