package bee2

import (
	"testing"
	"testing/quick"
)

var identityECB, checkECB = makeCryptoHelpers(nil,
	func(out, src []byte, key BeltKey, iv BeltIV, opt *CommonOpt) error {
		return ECBEncr(out, src, key, opt)
	},
	func(out, src []byte, key BeltKey, iv BeltIV, opt *CommonOpt) error {
		return ECBDecr(out, src, key, opt)
	},
)

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
		want := mustDecode(t, cases[i].want)

		checkECB(t, input, want, key, BeltIV{}, enc)
	}
}

func TestECBProp(t *testing.T) {
	encBuf := make([]byte, maxSize)

	f := func(input []byte, pass []byte) (ok bool) {
		if len(input) < 16 {
			return true
		}

		key := mustDerive(t, pass)
		enc := encBuf[:len(input)]

		identityECB(t, input, key, BeltIV{}, enc)

		return true
	}

	conf := conf(maxSize, 2)
	if err := quick.Check(f, &conf); err != nil {
		t.Fatal(err)
	}
}
