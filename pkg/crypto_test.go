package pkg

import (
	"fmt"
	"regexp"
	"testing"
)

func TestEncrypt(t *testing.T) {
	input := []byte("🙏🙏🙏")
	want := regexp.MustCompile("eb6ba8bde3821909b63e14764485530fd8e875a23834d41d6c100ac446828c7e")
	msg, err := Encrypt(input)
	res := fmt.Sprintf("%x", msg)
	if !want.MatchString(res) || err != nil {
		t.Fatalf("%q, %v, not %#q", res, err, want)
	}
}
