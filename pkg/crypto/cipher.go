package crypto

import (
	"github.com/head-gardener/passage/pkg/bee2/belt"

	"lukechampine.com/uint128"
)

type MAC = belt.MAC

type Cipher interface {
	Wrap(out []byte, crit []byte, open []byte, mac []byte) error
	Unwrap(out []byte, crit []byte, open []byte, mac belt.MAC) error
	Finalize()
}

type CHE struct {
	ivu belt.IV
	ivw belt.IV
	key belt.Key
}

func inc128(i *belt.IV) {
	uint128.FromBytes(i[:]).AddWrap64(1).PutBytes(i[:])
}

func InitCHE(pass []byte, salt []byte) (che *CHE, err error) {
	key, err := belt.KDF(pass, salt, nil)
	if err != nil {
		return
	}

	return &CHE{
		key: key,
		ivu: belt.IV{},
		ivw: belt.IV{},
	}, nil
}

func (che *CHE) Wrap(out []byte, crit []byte, open []byte, mac []byte) (err error) {
	m, err := belt.CHEWrap(out, crit, open, che.key, che.ivw, nil)
	if err != nil {
		return
	}

	copy(mac, m[:])
	inc128(&che.ivw)
	return nil
}

func (che *CHE) Unwrap(out []byte, crit []byte, open []byte, mac belt.MAC) (err error) {
	err = belt.CHEUnwrap(out, crit, open, mac, che.key, che.ivu, nil)
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
