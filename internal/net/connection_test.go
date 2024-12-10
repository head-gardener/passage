package net

import (
	"bytes"
	"net"
	"testing"
	"testing/quick"
	"time"

	"github.com/flynn/noise"

	"github.com/head-gardener/passage/internal/handshake"
	"github.com/head-gardener/passage/pkg/bee2/belt"
)

func (conn *Connection) mock(remote net.Conn, passx []byte, passr []byte) (err error) {
	var psk [32]byte
	err = belt.Hash(psk[:], passx, nil)
	if err != nil {
		return err
	}
	conn.tcp = remote
	conn.cx = noise.UnsafeNewCipherState(handshake.BignBeltSuite, psk, 0)
	err = belt.Hash(psk[:], passr, nil)
	conn.cr = noise.UnsafeNewCipherState(handshake.BignBeltSuite, psk, 0)
	return
}

func TestConnetionHandshake(t *testing.T) {
	nodeA := Connection{}
	nodeB := Connection{}

	pipeA, pipeB := net.Pipe()

	pipeA.SetDeadline(time.Now().Add(time.Second))
	pipeB.SetDeadline(time.Now().Add(time.Second))

	pass := []byte("pass")
	nodeA.tcp = pipeA
	nodeB.tcp = pipeB
	var psk [32]byte
	err := belt.Hash(psk[:], pass, nil)
	if err != nil {
		t.Fatal(err)
	}

	errs := make(chan error, 1)
	go func() {
		cx, cr, err := handshakeResponder(pipeB, psk[:])
		nodeB.cr = cr
		nodeB.cx = cx
		errs <- err
	}()
	cx, cr, err := handshakeInitiator(pipeA, psk[:])
	nodeA.cx = cx
	nodeA.cr = cr
	{
		err := <-errs
		if err != nil {
			t.Fatalf("handshake responder, %v", err)
		}
	}
	if err != nil {
		t.Fatalf("handshake initiator, %v", err)
	}

	b := []byte("handshake-test")

	output := make([]byte, 100)
	input := make([]byte, 100)
	copy(input, b)

	errs = make(chan error, 1)
	body := len(b)

	go func() {
		_, err := nodeB.Read(output)
		errs <- err
	}()

	_, err = nodeA.Write(input[:body])
	{
		err := <-errs
		if err != nil {
			t.Fatalf("error reading, %v", err)
		}
	}
	if err != nil {
		t.Fatalf("error writing, %v", err)
	}

	if !bytes.Equal(b, output[:body]) {
		t.Errorf("output mismatch: %x, %x", b, output[:body])
		return
	}
}

func TestConnetionProp(t *testing.T) {
	nodeA := Connection{}
	nodeB := Connection{}

	pipeA, pipeB := net.Pipe()

	pipeA.SetDeadline(time.Now().Add(time.Second))
	pipeB.SetDeadline(time.Now().Add(time.Second))

	passA := []byte("passA")
	passB := []byte("passB")
	nodeA.mock(pipeA, passA, passB)
	nodeB.mock(pipeB, passA, passB)

	output := make([]byte, 100)
	input := make([]byte, 100)

	f := func(b []byte) (ok bool) {
		if len(b) == 0 {
			return true
		}
		ok = false
		copy(input, b)
		errs := make(chan error, 1)
		body := len(b)

		go func() {
			_, err := nodeB.Read(output)
			errs <- err
		}()

		_, err := nodeA.Write(input[:body])
		{
			err := <-errs
			if err != nil {
				t.Fatalf("error reading, %v", err)
			}
		}
		if err != nil {
			t.Fatalf("error writing, %v", err)
		}

		if !bytes.Equal(b, output[:body]) {
			t.Errorf("output mismatch: %x, %x", b, output[:body])
			return
		}

		ok = true
		return
	}

	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}
