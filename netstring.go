/*
D. J. Bernstein's netstrings for Go.
*/
package netstring

import (
	"errors"
	"io"
)

var Incomplete = errors.New("The netstring is incomplete")
var Garbled = errors.New("The netstring was not correctly formatted and could not be read")

type Netstring struct {
	buffer   []byte
	complete bool
}

// Construct a netstring wrapping a byte slice, for output.
// n.IsComplete() will return true.
func From(buf []byte) *Netstring {
	return &Netstring{
		buffer:   buf,
		complete: true,
	}
}

// Construct an empty netstring, for input.
// n.IsComplete() will return false.
func ForReading() *Netstring {
	return &Netstring{buffer: nil, complete: false}
}

// Returns true if the number of bytes advertized in the netstring's length have been read into its buffer.
// Operations that require the netstring's contents will not be available until it is complete.
func (n *Netstring) IsComplete() bool { return n.complete }

// Returns the advertized length of the netstring.
// If n.IsComplete() then this is the length of the data in the buffer, too.
// If this value is negative, in means that no length has been read yet.
func (n *Netstring) Length() int {
	if n.buffer == nil {
		return -1
	}
	return cap(n.buffer)
}

// Returns the bytes in the netstring if it is complete, otherwise returns Incomplete as an error.
func (n *Netstring) Bytes() ([]byte, error) {
	if n.complete {
		return n.buffer, nil
	}
	return nil, Incomplete
}

// Read a netstring from input.
// Returns any errors from input except io.EOF.
// Returns Garbled if the input was not a valid netstring.
// Returns Incomplete if the input was shorter than a full netstring.
// To resume reading where you left off, call Readfrom(input) again.
// Calling Readfrom(input) on a complete netstring does nothing.
func (n *Netstring) ReadFrom(input io.Reader) error {
	var err error
	if n.buffer == nil {
		var length int
		length, err = n.readLength(input)
		if err != nil {
			return err
		}
		err = n.readColon(input)
		if err != nil {
			return err
		}
		n.buffer = make([]byte, 0, length) // capacity stores the length
	}
	if len(n.buffer) < cap(n.buffer) {
		// slice n.buffer to the part between len and cap
		dest := n.buffer[len(n.buffer):cap(n.buffer)]
		count, err := input.Read(dest)

		// slice n.buffer to add on count bytes
		if count > 0 {
			n.buffer = n.buffer[:len(n.buffer)+count]
		}

		switch {
		case err == io.EOF: // we still expect to read a comma, so EOF here is always incomplete
			return Incomplete
		case len(n.buffer) < cap(n.buffer):
			return Incomplete
		case err != nil:
			return err
		}
	}
	if !n.complete {
		err = n.readComma(input)
		if err != nil {
			return err
		}
		n.complete = true
	}
	return nil
}

func (n *Netstring) readLength(input io.Reader) (int, error) {
	//TODO
	return 0, nil
}

func (n *Netstring) readColon(input io.Reader) error {
	//TODO
	return nil
}

func (n *Netstring) readComma(input io.Reader) error {
	//TODO
	return nil
}
