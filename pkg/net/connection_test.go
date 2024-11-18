package net

import (
	"bytes"
	"net"
	"testing"
	"testing/quick"
	"time"
)

func TestConnetionHandshake(t *testing.T) {
	nodeA := Connection{}
	nodeB := Connection{}

	pipeA, pipeB := net.Pipe()

	pipeA.SetDeadline(time.Now().Add(time.Second))
	pipeB.SetDeadline(time.Now().Add(time.Second))

	pass := []byte("pass")
	nodeA.mock(pipeA, pass)
	nodeB.mock(pipeB, pass)

	errs := make(chan error, 1)
	go func() {
		cipher, err := handshakeResponder(pipeB, pass)
		nodeB.cipher = cipher
		errs <- err
	}()
	cipher, err := handshakeInitiator(pipeA, pass)
	nodeA.cipher = cipher
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

	pass := []byte("pass")
	nodeA.mock(pipeA, pass)
	nodeB.mock(pipeB, pass)

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
