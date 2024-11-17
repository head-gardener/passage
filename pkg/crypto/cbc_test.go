package crypto

import (
	"bytes"
	"testing"
	"testing/quick"
)

var propCBCIdentity = makeCryptoIdentity(nil, CBCEncr, CBCDecr)

func TestCBC(t *testing.T) {
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
			want:  "10116EFAE6AD58EE14852E11DA1B8A745CF2480E8D03F1C19492E53ED3A70F60657C1EE8C0E0AE5B58388BF8A68E3309",
		},
		{
			input: "B194BAC80A08F53B366D008E584A5DE48504FA9D1BB6C7AC252E72C202FDCE0D5BE3D612",
			key:   "E9DEE72C8F0C0FA62DDB49F46F73964706075316ED247A3739CBA38303A98BF6",
			iv:    "BE32971343FC9A48A02A885F194B09A1",
			want:  "10116EFAE6AD58EE14852E11DA1B8A746A9BBADCAF73F968F875DEDC0A44F6B15CF2480E",
		},
		{
			input: "730894D6158E17CC1600185A8F411CAB0471FF85C83792398D8924EBD57D03DB95B97A9B7907E4B020960455E46176F8",
			key:   "92BD9B1CE5D141015445FBC95E4D0EF2682080AA227D642F2687F93490405511",
			iv:    "7ECDA4D01544AF8CA58450BF66D2E88A",
			want:  "E12BDC1AE28257EC703FCCF095EE8DF1C1AB76389FE678CAF7C6F860D5BB9C4FF33C657B637C306ADD4EA7799EB23D31",
		},
		{
			input: "730894D6158E17CC1600185A8F411CABB6AB7AF8541CF85755B8EA27239F08D2166646E4",
			key:   "92BD9B1CE5D141015445FBC95E4D0EF2682080AA227D642F2687F93490405511",
			iv:    "7ECDA4D01544AF8CA58450BF66D2E88A",
			want:  "E12BDC1AE28257EC703FCCF095EE8DF1C1AB76389FE678CAF7C6F860D5BB9C4FF33C657B",
		},
	}
	for i := range cases {
		input := mustDecode(t, cases[i].input)
		key := mustBeKey(t, cases[i].key)
		iv := mustBeIV(t, cases[i].iv)
		enc := make([]byte, len(input))
		dec := make([]byte, len(input))

		propCBCIdentity(t, input, key, iv, enc, dec)

		want := mustDecode(t, cases[i].want)
		if !bytes.Equal(enc, want) {
			t.Fatalf("no match:\n%x + %x ->\n%x, not\n%x", input, key, enc, want)
		}
	}
}

func TestCBCProp(t *testing.T) {
	encBuf := make([]byte, maxSize)
	decBuf := make([]byte, maxSize)

	f := func(input []byte, pass []byte) (ok bool) {
		if len(input) < 16 {
			return true
		}

		key := mustDerive(t, pass)
		iv := mustContainIV(t, encBuf)
		enc := encBuf[:len(input)]
		dec := decBuf[:len(input)]

		propCBCIdentity(t, input, key, iv, enc, dec)

		return true
	}

	conf := conf(maxSize, 3)
	if err := quick.Check(f, &conf); err != nil {
		t.Fatal(err)
	}
}
