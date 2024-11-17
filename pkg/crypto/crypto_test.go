package crypto

import (
	"bytes"
	"encoding/hex"
	"math/rand"
	"reflect"
	"testing"
	"testing/quick"
)

const maxSize = 100000

func conf(maxSize int, argn int) quick.Config {
	buf := make([]byte, maxSize * argn)
	return quick.Config{
		Values: func(args []reflect.Value, rand *rand.Rand) {
			for i := range args {
				offset := maxSize * i
				n := rand.Intn(maxSize) + 1
				slice := buf[offset : offset+n]
				rand.Read(slice)
				args[i] = reflect.ValueOf(slice)
			}
		},
	}
}

func mustDecode(t *testing.T, str string) (res []byte) {
	res, err := hex.DecodeString(str)
	if err != nil {
		t.Fatal(err)
	}
	return
}

func mustBeKey(t *testing.T, str string) (key BeltKey) {
	keyBuf := mustDecode(t, str)
	if len(keyBuf) != 32 {
		t.Fatalf("invalid length %d for belt key", len(keyBuf))
	}
	copy(key[:], keyBuf)
	return key
}

func mustDerive(t *testing.T, str []byte) (key BeltKey) {
	key, err := KDF(str, []byte{0}, &KDFOpt{iter: 1})
	if err != nil {
		t.Fatal(err)
	}
	return
}

func propECBIdentity(t *testing.T, input []byte, key BeltKey, enc []byte, dec []byte) {
	if err := ECBEncr(enc, input, key, nil); err != nil {
		t.Fatal(err)
	}
	if err := ECBDecr(dec, enc, key, nil); err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(dec, input) {
		t.Fatalf("no decryption:\n%x + %x ->\n%x, not\n%x", enc, key, dec, input)
	}

	// in-place enc/dec
	if err := ECBEncr(dec, dec, key, nil); err != nil {
		t.Fatal(err)
	}
	if err := ECBDecr(dec, dec, key, nil); err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(dec, input) {
		t.Fatalf("in-place no decryption:\n%x + %x ->\n%x, not\n%x", enc, key, dec, input)
	}
}

func TestECB(t *testing.T) {
	cases := []struct {
		input string
		key   string
		want  string
	}{
		{
			input: "B194BAC80A08F53B366D008E584A5DE48504FA9D1BB6C7AC252E72C202FDCE0D5BE3D61217B96181FE6786AD716B890B",
			key:   "E9DEE72C8F0C0FA62DDB49F46F73964706075316ED247A3739CBA38303A98BF6",
			want:  "69CCA1C93557C9E3D66BC3E0FA88FA6E5F23102EF109710775017F73806DA9DC46FB2ED2CE771F26DCB5E5D1569F9AB0",
		},
		{
			input: "B194BAC80A08F53B366D008E584A5DE48504FA9D1BB6C7AC252E72C202FDCE0D5BE3D61217B96181FE6786AD716B89",
			key:   "E9DEE72C8F0C0FA62DDB49F46F73964706075316ED247A3739CBA38303A98BF6",
			want:  "69CCA1C93557C9E3D66BC3E0FA88FA6E36F00CFED6D1CA1498C12798F4BEB2075F23102EF109710775017F73806DA9",
		},
		{
			input: "0DC5300600CAB840B38448E5E993F421E55A239F2AB5C5D5FDB6E81B40938E2A54120CA3E6E19C7AD750FC3531DAEAB7",
			want:  "E12BDC1AE28257EC703FCCF095EE8DF1C1AB76389FE678CAF7C6F860D5BB9C4FF33C657B637C306ADD4EA7799EB23D31",
			key:   "92BD9B1CE5D141015445FBC95E4D0EF2682080AA227D642F2687F93490405511",
		},
		{
			input: "0DC5300600CAB840B38448E5E993F4215780A6E2B69EAFBB258726D7B6718523E55A239F",
			want:  "E12BDC1AE28257EC703FCCF095EE8DF1C1AB76389FE678CAF7C6F860D5BB9C4FF33C657B",
			key:   "92BD9B1CE5D141015445FBC95E4D0EF2682080AA227D642F2687F93490405511",
		},
	}

	for i := range cases {
		input := mustDecode(t, cases[i].input)
		key := mustBeKey(t, cases[i].key)
		enc := make([]byte, len(input))
		dec := make([]byte, len(input))

		propECBIdentity(t, input, key, enc, dec)

		want := mustDecode(t, cases[i].want)
		if !bytes.Equal(enc, want) {
			t.Fatalf("no match:\n%x + %x ->\n%x, not\n%x", input, key, enc, want)
		}
	}
}

func TestECBProp(t *testing.T) {
	encBuf := make([]byte, maxSize)
	decBuf := make([]byte, maxSize)

	f := func(input []byte, pass []byte) (ok bool) {
		if len(input) < 16 {
			return true
		}

		key := mustDerive(t, pass)
		enc := encBuf[:len(input)]
		dec := decBuf[:len(input)]

		propECBIdentity(t, input, key, enc, dec)

		return true
	}

	conf := conf(maxSize, 2)
	if err := quick.Check(f, &conf); err != nil {
		t.Fatal(err)
	}
}

func TestKDF(t *testing.T) {
	cases := []struct {
		pass string
		salt string
		want string
	}{
		{
			pass: "42313934424143383041303846353342",
			salt: "BE32971343FC9A48",
			want: "3D331BBBB1FBBB40E4BF22F6CB9A689EF13A77DC09ECF93291BFE42439A72E7D",
		},
	}

	for i := range cases {
		pass := mustDecode(t, cases[i].pass)
		salt := mustDecode(t, cases[i].salt)

		key, err := KDF(pass, salt, nil)
		if err != nil {
			t.Fatal(err)
		}

		want := mustDecode(t, cases[i].want)
		keyBuf := make([]byte, 32)
		copy(keyBuf, key[:])
		if !bytes.Equal(keyBuf, want) {
			t.Fatalf("no match:\n%x + %x ->\n%x, not\n%x", pass, salt, key, want)
		}
	}
}

func TestKDFProp(t *testing.T) {
	encA := make([]byte, 16)
	encB := make([]byte, 16)

	f := func(pass []byte) (ok bool) {
		key := mustDerive(t, pass)

		if err := ECBEncr(encB, encA, key, nil); err != nil {
			t.Fatal(err)
		}
		if bytes.Equal(encA, encB) {
			t.Fatalf("key %x didn't change input %x", key, encA)
		}
		copy(encA, encB)

		return true
	}

	conf := conf(maxSize, 1)
	if err := quick.Check(f, &conf); err != nil {
		t.Fatal(err)
	}
}

func TestHMAC(t *testing.T) {
	cases := []struct {
		input string
		key   string
		want  string
	}{
		{
			input: "BE32971343FC9A48A02A885F194B09A17ECDA4D01544AF8CA58450BF66D2E88A",
			key:   "E9DEE72C8F0C0FA62DDB49F46F73964706075316ED247A3739CBA38303",
			want:  "D4828E6312B08BB83C9FA6535A4635549E411FD11C0D8289359A1130E930676B",
		},
		{
			input: "BE32971343FC9A48A02A885F194B09A17ECDA4D01544AF8CA58450BF66D2E88A",
			key:   "E9DEE72C8F0C0FA62DDB49F46F73964706075316ED247A3739CBA38303A98BF6",
			want:  "41FFE8645AEC0612E952D2CDF8DD508F3E4A1D9B53F6A1DB293B19FE76B1879F",
		},
		{
			input: "BE32971343FC9A48A02A885F194B09A17ECDA4D01544AF8CA58450BF66D2E88A",
			key:   "E9DEE72C8F0C0FA62DDB49F46F73964706075316ED247A3739CBA38303A98BF692BD9B1CE5D141015445",
			want:  "7D01B84D2315C332277B3653D7EC64707EBA7CDFF7FF70077B1DECBD68F2A144",
		},
	}

	for i := range cases {
		out := make([]byte, 32)
		input := mustDecode(t, cases[i].input)
		key := mustDecode(t, cases[i].key)

		if err := HMAC(out, input, key, nil); err != nil {
			t.Fatal(err)
		}

		want := mustDecode(t, cases[i].want)
		if !bytes.Equal(out, want) {
			t.Fatalf("no match:\n%x + %x ->\n%x, not\n%x", input, key, out, want)
		}
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

		want := mustDecode(t, cases[i].want)
		if !bytes.Equal(out, want) {
			t.Fatalf("no match:\n%x ->\n%x, not\n%x", input, out, want)
		}
		if !bytes.Equal(out2[:32], want) {
			t.Fatalf("in-place no match:\n%x ->\n%x, not\n%x", input, out2, want)
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
