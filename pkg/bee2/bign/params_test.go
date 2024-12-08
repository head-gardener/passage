package bign

import (
	"testing"
)

func TestStandardParams(t *testing.T) {
	p, err := StandardParams(L128)
	if err != nil {
		t.Fatal(err)
	}
	if err := p.Validate(); err != nil {
		t.Fatal(err)
	}
	p, err = StandardParams(L192)
	if err != nil {
		t.Fatal(err)
	}
	if err := p.Validate(); err != nil {
		t.Fatal(err)
	}
	p, err = StandardParams(L256)
	if err != nil {
		t.Fatal(err)
	}
	if err := p.Validate(); err != nil {
		t.Fatal(err)
	}
}
