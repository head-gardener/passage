package handshake

import (
	"crypto/rand"
	"testing"

	"github.com/head-gardener/passage/pkg/bee2/belt"

	. "gopkg.in/check.v1"
)

func TestNoiseSuite(t *testing.T) { TestingT(t) }

type NoiseSuite struct{}

var _ = Suite(&NoiseSuite{})

func (NoiseSuite) TestNNpsk0Roundtrip(c *C) {
	psk := make([]byte, 32)
	rand.Read(psk)

	hsI, err := Init(true, psk)
	hsR, err := Init(false, psk)
	overhead := belt.NewCHE(belt.Key(psk)).Overhead()

	// -> e
	msg, _, _, _ := hsI.WriteMessage(nil, nil)
	c.Assert(msg, HasLen, Bign128.DHLen()+overhead)
	res, _, _, err := hsR.ReadMessage(nil, msg)
	c.Assert(err, IsNil)
	c.Assert(res, HasLen, 0)

	// <- e, dhee
	msg, csR0, csR1, _ := hsR.WriteMessage(nil, nil)
	c.Assert(msg, HasLen, Bign128.DHLen()+overhead)
	res, csI0, csI1, err := hsI.ReadMessage(nil, msg)
	c.Assert(err, IsNil)
	c.Assert(res, HasLen, 0)

	// transport I -> R
	msg, err = csI0.Encrypt(nil, nil, []byte("foo"))
	c.Assert(err, IsNil)
	res, err = csR0.Decrypt(nil, nil, msg)
	c.Assert(err, IsNil)
	c.Assert(string(res), Equals, "foo")

	// transport R -> I
	msg, err = csR1.Encrypt(nil, nil, []byte("bar"))
	c.Assert(err, IsNil)
	res, err = csI1.Decrypt(nil, nil, msg)
	c.Assert(err, IsNil)
	c.Assert(string(res), Equals, "bar")
}
