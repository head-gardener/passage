package belt

import (
	"bytes"
	"testing"
	"testing/quick"
)

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
