package belt

import (
	"testing"
	"testing/quick"
)

var identityCFB, checkCFB = makeCryptoHelpers(nil, CFBEncr, CFBDecr)

func TestCFB(t *testing.T) {
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
			want:  "C31E490A90EFA374626CC99E4B7B8540A6E48685464A5A06849C9CA769A1B0AE55C2CC5939303EC832DD2FE16C8E5A1B",
		},
		{
			input: "FA9D107A86F375EE65CD1DB881224BD016AFF814938ED39B3361ABB0BF0851B652244EB06842DD4C94AA4500774E40BB",
			key:   "92BD9B1CE5D141015445FBC95E4D0EF2682080AA227D642F2687F93490405511",
			iv:    "7ECDA4D01544AF8CA58450BF66D2E88A",
			want:  "E12BDC1AE28257EC703FCCF095EE8DF1C1AB76389FE678CAF7C6F860D5BB9C4FF33C657B637C306ADD4EA7799EB23D31",
		},
	}
	for i := range cases {
		input := mustDecode(t, cases[i].input)
		key := mustBeKey(t, cases[i].key)
		iv := mustBeIV(t, cases[i].iv)
		enc := make([]byte, len(input))
		want := mustDecode(t, cases[i].want)

		checkCFB(t, input, want, key, iv, enc)
	}
}

func TestCFBProp(t *testing.T) {
	encBuf := make([]byte, maxSize)

	f := func(input []byte, pass []byte) (ok bool) {
		key := mustDerive(t, pass)
		iv := mustContainIV(t, encBuf)
		enc := encBuf[:len(input)]

		identityCFB(t, input, key, iv, enc)

		return true
	}

	conf := conf(maxSize, 3)
	if err := quick.Check(f, &conf); err != nil {
		t.Fatal(err)
	}
}
