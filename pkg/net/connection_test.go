package net

import (
	"bytes"
	"net"
	"testing"
	"testing/quick"
	"time"
)

func TestConnetion(t *testing.T) {
	nodeA := Connection{}
	nodeB := Connection{}

	pipeA, pipeB := net.Pipe()

	pipeA.SetDeadline(time.Now().Add(time.Second))
	pipeB.SetDeadline(time.Now().Add(time.Second))

	nodeA.mock(pipeA)
	nodeB.mock(pipeB)

	f := func(b []byte) (ok bool) {
		ok = false
		output := make([]byte, len(b))
		errs := make(chan error, 1)

		go func() {
			_, err := nodeB.Read(output)
			if err != nil {
				errs <- err
			}
		}()

		_, err := nodeA.Write(b)
		if err != nil {
			t.Fatalf("error writing, %v", err)
		}

		select {
		case err = <-errs:
			t.Fatalf("error reading, %v", err)
		default:
		}

		if !bytes.Equal(b, output) {
			return
		}

		ok = true
		return
	}

	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}
