package bee2

// Produced by KDF, consumed by everything else. Same key shouldn't be used in different
// algorithms. A 32 byte, 256 bit slice.
type BeltKey [32]byte

// Must be unique for every session using a single key. A 16 byte, 128 bit slice.
type BeltIV [16]byte

// Produced and consumed by AEAD functions. A 8 byte, 64 bit slice.
type BeltMAC [8]byte

type CommonOpt struct {
	srcLen int
}

type AEADOpt struct {
	critLen int
	openLen int
}
