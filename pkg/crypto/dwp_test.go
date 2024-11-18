package crypto

import (
	"testing"
	"testing/quick"
)

var identityDWP, checkDWP = makeCryptoHelpers(nil,
	func(out *dataAEAD, src *dataAEAD, key BeltKey, iv BeltIV, opt *AEADOpt) error {
		mac, err := DWPWrap(out.crit, src.crit, src.open, key, iv, opt)
		out.mac = mac
		out.open = src.open
		return err
	},
	func(out *dataAEAD, src *dataAEAD, key BeltKey, iv BeltIV, opt *AEADOpt) error {
		return DWPUnwrap(out.crit, src.crit, src.open, src.mac, key, iv, opt)
	},
)

func TestDWP(t *testing.T) {
	cases := []struct {
		crit string
		open string
		key  string
		iv   string
		mac  string
		want string
	}{
		{
			open: "8504FA9D1BB6C7AC252E72C202FDCE0D5BE3D61217B96181FE6786AD716B890B",
			key:  "E9DEE72C8F0C0FA62DDB49F46F73964706075316ED247A3739CBA38303A98BF6",
			iv:   "BE32971343FC9A48A02A885F194B09A1",
			crit: "B194BAC80A08F53B366D008E584A5DE4",
			want: "52C9AF96FF50F64435FC43DEF56BD797",
			mac:  "3B2E0AEB2B91854B",
		},
		{
			open: "C1AB76389FE678CAF7C6F860D5BB9C4FF33C657B637C306ADD4EA7799EB23D31",
			key:  "92BD9B1CE5D141015445FBC95E4D0EF2682080AA227D642F2687F93490405511",
			iv:   "7ECDA4D01544AF8CA58450BF66D2E88A",
			crit: "DF181ED008A20F43DCBBB93650DAD34B",
			want: "E12BDC1AE28257EC703FCCF095EE8DF1",
			mac:  "6A2C2C94C4150DC0",
		},
	}
	for i := range cases {
		key := mustBeKey(t, cases[i].key)
		iv := mustBeIV(t, cases[i].iv)
		input := dataAEAD{
			crit: mustDecode(t, cases[i].crit),
			open: mustDecode(t, cases[i].open),
			mac:  BeltMAC{},
		}
		enc := dataAEAD{
			crit: make([]byte, len(input.crit)),
			open: make([]byte, len(input.open)),
			mac:  BeltMAC{},
		}
		want := dataAEAD{
			crit: mustDecode(t, cases[i].want),
			open: mustDecode(t, cases[i].open),
			mac:  mustBeMAC(t, cases[i].mac),
		}

		checkDWP(t, &input, &want, key, iv, &enc)
	}
}

func TestDWPProp(t *testing.T) {
	encOpenBuf := make([]byte, maxSize)
	encCritBuf := make([]byte, maxSize)

	f := func(crit []byte, open []byte, pass []byte) (ok bool) {
		key := mustDerive(t, pass)
		iv := mustContainIV(t, encOpenBuf)
		input := dataAEAD{
			crit: crit,
			open: open,
			mac:  BeltMAC{},
		}
		enc := dataAEAD{
			crit: encCritBuf[:len(input.crit)],
			open: encOpenBuf[:len(input.crit)],
			mac:  BeltMAC{},
		}

		identityDWP(t, &input, key, iv, &enc)

		return true
	}

	conf := conf(maxSize, 3)
	if err := quick.Check(f, &conf); err != nil {
		t.Fatal(err)
	}
}
