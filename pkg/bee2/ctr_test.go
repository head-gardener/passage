package bee2

import (
	"testing"
	"testing/quick"
)

var identityCTR, checkCTR = makeCryptoHelpers(nil, CTR, CTR)

func TestCTR(t *testing.T) {
	cases := []struct {
		input string
		key   string
		iv    string
		want  string
	}{
		{
			input: "B194BAC80A08F53B366D008E584A5DE48504FA9D1BB6C7AC252E72C202FDCE0D5BE3D61217B96181FE6786AD716B890B",
			key:   "E9DEE72C8F0C0FA62DDB49F46F73964706075316ED247A3739CBA38303A98BF6",
			iv:    "BE32971343FC9A48A02A885F194B09A1",
			want:  "52C9AF96FF50F64435FC43DEF56BD797D5B5B1FF79FB41257AB9CDF6E63E81F8F00341473EAE409833622DE05213773A",
		},
	}
	for i := range cases {
		input := mustDecode(t, cases[i].input)
		key := mustBeKey(t, cases[i].key)
		iv := mustBeIV(t, cases[i].iv)
		enc := make([]byte, len(input))
		want := mustDecode(t, cases[i].want)

		checkCTR(t, input, want, key, iv, enc)
	}
}

func TestCTRProp(t *testing.T) {
	encBuf := make([]byte, maxSize)

	f := func(input []byte, pass []byte) (ok bool) {
		key := mustDerive(t, pass)
		iv := mustContainIV(t, encBuf)
		enc := encBuf[:len(input)]

		identityCTR(t, input, key, iv, enc)

		return true
	}

	conf := conf(maxSize, 3)
	if err := quick.Check(f, &conf); err != nil {
		t.Fatal(err)
	}
}
