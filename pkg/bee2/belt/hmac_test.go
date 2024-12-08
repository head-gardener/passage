package belt

import (
	"bytes"
	"testing"
)

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
