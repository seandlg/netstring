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

func TestIncomplete(t *testing.T) {
	n := ForReading()
	input := bytes.NewBufferString("1")
	err := n.ReadFrom(input)
	if err != Incomplete {
		t.Fatal(err)
	}
	input.WriteString("1:hello")
	err = n.ReadFrom(input)
	if err != Incomplete {
		t.Fatal(err)
	}
	input.WriteString(" world")
	err = n.ReadFrom(input)
	if err != Incomplete {
		t.Fatal(err)
	}
	input.WriteString(",")
	err = n.ReadFrom(input)
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
