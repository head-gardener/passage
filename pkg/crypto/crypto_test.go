package crypto

import (
	"bytes"
	"encoding/hex"
	"testing"
	"testing/quick"
)

func TestHash(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{
			input: "B194BAC80A08F53B366D008E58",
			want:  "ABEF9725D4C5A83597A367D14494CC2542F20F659DDFECC961A3EC550CBA8C75",
		},
		{
			input: "B194BAC80A08F53B366D008E584A5DE48504FA9D1BB6C7AC252E72C202FDCE0D",
			want:  "749E4C3653AECE5E48DB4761227742EB6DBE13F4A80F7BEFF1A9CF8D10EE7786",
		},
		{
			input: "B194BAC80A08F53B366D008E584A5DE48504FA9D1BB6C7AC252E72C202FDCE0D5BE3D61217B96181FE6786AD716B890B",
			want:  "9D02EE446FB6A29FE5C982D4B13AF9D3E90861BC4CEF27CF306BFB0B174A154A",
		},
	}

	for i := range cases {
		out := make([]byte, 32)
		input, err := hex.DecodeString(cases[i].input)
		if err != nil {
			t.Fatal(err)
		}

		want, err := hex.DecodeString(cases[i].want)
		if err != nil {
			t.Fatal(err)
		}

		if err := Hash(input, out); err != nil {
			t.Fatal(err)
		}

		if !bytes.Equal(out, want) {
			t.Fatalf("no match: %x -> %x, not %x", input, out, want)
		}
	}
}

func TestHashProp(t *testing.T) {
	first := make([]byte, 32)
	second := make([]byte, 32)

	f := func(b []byte) (ok bool) {
		if len(b) == 0 {
			return true
		}
		if err := Hash(b, first); err != nil {
			t.Fatal(err)
		}
		if err := Hash(b, second); err != nil {
			t.Fatal(err)
		}
		return bytes.Equal(first, second)
	}

	if err := quick.Check(f, nil); err != nil {
		t.Fatal(err)
	}
}
