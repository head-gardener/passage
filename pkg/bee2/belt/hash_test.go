package belt

import (
	"bytes"
	"testing"
	"testing/quick"
)

func TestHashInterface(t *testing.T) {
	h := HashInit()
	input := mustDecode(t, "B194BAC80A08F53B366D008E58")
	want := mustDecode(t, "ABEF9725D4C5A83597A367D14494CC2542F20F659DDFECC961A3EC550CBA8C75")

	h.Write(input)
	out1 := h.Sum(nil)
	h.Reset()
	if len(out1) != h.Size() {
		t.Fatalf("Size() is incorrect")
	}

	h.Write(input)
	out2 := h.Sum(nil)
	h.Reset()
	if !bytes.Equal(out1, out2) {
		t.Fatalf("Reset() doesn't reset state")
	}

	h.Write(input)
	if !bytes.Equal(h.Sum(nil), h.Sum(nil)) {
		t.Fatalf("Sum() modifies state")
	}
	h.Reset()

	h.Write(input)
	out1 = h.Sum(input)
	h.Reset()
	if !bytes.Equal(out1, append(input, want...)) {
		t.Fatalf("Sum() doesn't append")
	}

	h.Write(input)
	h.Write(input)
	out1 = h.Sum(input)
	h.Reset()
	h.Write(append(input, input...))
	out2 = h.Sum(input)
	h.Reset()
	if !bytes.Equal(out1, out2) {
		t.Fatalf("double Write() is not equivalent to double input")
	}
}

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
		input := mustDecode(t, cases[i].input)
		if err := Hash(out, input, nil); err != nil {
			t.Fatal(err)
		}

		out2 := make([]byte, max(len(input), 32))
		copy(out2, input)
		if err := Hash(out2, out2, &CommonOpt{srcLen: len(input)}); err != nil {
			t.Fatal(err)
		}

		h := HashInit()
		h.Write(input)
		out3 := h.Sum(nil)

		want := mustDecode(t, cases[i].want)
		if !bytes.Equal(out, want) {
			t.Fatalf("no match:\n%x ->\n%x, not\n%x", input, out, want)
		}
		if !bytes.Equal(out2[:32], want) {
			t.Fatalf("in-place no match:\n%x ->\n%x, not\n%x", input, out2, want)
		}
		if !bytes.Equal(out3[:32], want) {
			t.Fatalf("interface no match:\n%x ->\n%x, not\n%x", input, out2, want)
		}
	}
}

func TestHashProp(t *testing.T) {
	first := make([]byte, 32)
	second := make([]byte, 32)
	third := make([]byte, maxSize)

	f := func(b []byte) (ok bool) {
		if len(b) == 0 {
			return true
		}
		if err := Hash(first, b, nil); err != nil {
			t.Fatal(err)
		}
		if err := Hash(second, b, nil); err != nil {
			t.Fatal(err)
		}
		copy(third, b)
		if err := Hash(third, third[:len(b)], nil); err != nil {
			t.Fatal(err)
		}
		return bytes.Equal(first, second) && bytes.Equal(second, third[:32])
	}

	conf := conf(maxSize, 1)
	if err := quick.Check(f, &conf); err != nil {
		t.Fatal(err)
	}
}
