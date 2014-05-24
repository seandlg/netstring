/*
D. J. Bernstein's netstrings for Go.
*/
package netstring

import (
	"bytes"
	"errors"
)

var Incomplete = errors.New("The netstring is incomplete")

type Netstring struct {
	length int64
	buffer *bytes.Buffer
}

func From(buf []byte) *Netstring {
	return &Netstring{
		length: int64(len(buf)),
		buffer: bytes.NewBuffer(buf),
	}
}

func (n *Netstring) IsComplete() bool {
	if n.length < 0 || n.buffer == nil {
		return false
	}
	return n.length == int64(n.buffer.Len())
}

// Returns the advertized length of the netstring. If n.IsComplete() then this is the length of the data in the buffr, too.
func (n *Netstring) Length() int64 { return n.length }

func (n *Netstring) Bytes() []byte {
	return n.buffer.Bytes()
}
