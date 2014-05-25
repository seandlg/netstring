package netstring

import (
	"bytes"
	"testing"
)

func TestFromAndBytes(t *testing.T) {
	n := From([]byte("hello world"))
	out, err := n.Bytes()
	if err != nil {
		t.Fatal(err)
	}
	if "hello world" != string(out) {
		t.Fatal(out)
	}
}

func TestReading(t *testing.T) {
	n := ForReading()
	input := bytes.NewBufferString("11:hello world,")
	t.Log("length is: ", input.Len())
	err := n.ReadFrom(input)
	if err != nil {
		t.Fatal(err)
	}
	out, err := n.Bytes()
	if err != nil {
		t.Fatal(err)
	}
	if "hello world" != string(out) {
		t.Fatal(out)
	}
}
