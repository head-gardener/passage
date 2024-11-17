package crypto

import (
	"bytes"
	"encoding/hex"
	"math/rand"
	"reflect"
	"testing"
	"testing/quick"
)

func conf(maxSize int, buf []byte) quick.Config {
	return quick.Config{
		Values: func(args []reflect.Value, rand *rand.Rand) {
			for i := range args {
				n := rand.Intn(maxSize) + 1
				rand.Read(buf[:n])
				args[i] = reflect.ValueOf(buf[:n])
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

func TestHashInPlaceProp(t *testing.T) {
	first := make([]byte, 32)
	second := make([]byte, 100000)

	f := func(b []byte) (ok bool) {
		if len(b) == 0 {
			return true
		}
		if err := Hash(first, b, nil); err != nil {
			t.Fatal(err)
		}
		copy(second, b)
		if err := Hash(second, second, &CommonOpt{srcLen: len(b)}); err != nil {
			t.Fatal(err)
		}
		return bytes.Equal(first, second[:32])
	}

	buf := make([]byte, 100000)
	conf := conf(100000, buf)
	if err := quick.Check(f, &conf); err != nil {
		t.Fatal(err)
	}
}

func TestHashProp(t *testing.T) {
	first := make([]byte, 32)
	second := make([]byte, 32)

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
		return bytes.Equal(first, second)
	}

	buf := make([]byte, 100000)
	conf := conf(100000, buf)
	if err := quick.Check(f, &conf); err != nil {
		t.Fatal(err)
	}
}
