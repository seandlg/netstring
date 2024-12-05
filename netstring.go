/*
D. J. Bernstein's netstrings for Go.
*/
package netstring

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"math"
	"strconv"
)

// Convenience function, create a netstring representation of a []byte
func MarshalFrom(buf []byte) []byte {
	var out bytes.Buffer
	out.WriteString(strconv.Itoa(len(buf)))
	out.WriteRune(':')
	out.Write(buf)
	out.WriteRune(',')
	return out.Bytes()
}

var (
	Incomplete = errors.New("The netstring is incomplete")
	Garbled    = errors.New("The netstring was not correctly formatted and could not be read")
	TooLarge   = errors.New("The netstring length is too large")
)

type Netstring struct {
	maxLength int
	digits    []byte
	buffer    []byte
	complete  bool
}

// Construct a netstring wrapping a byte slice, for output.
// n.IsComplete() will return true.
func From(buf []byte) *Netstring {
	return &Netstring{
		digits:    make([]byte, 0, 10),
		buffer:    buf,
		complete:  true,
		maxLength: math.MaxInt,
	}
}

// Construct an empty netstring, for input.
// n.IsComplete() will return false.
func ForReading() *Netstring {
	return ForReadingWithMaxLength(math.MaxInt)
}

func ForReadingWithMaxLength(maxLength int) *Netstring {
	return &Netstring{digits: make([]byte, 0, 10), buffer: nil, complete: false, maxLength: maxLength}
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

// Formats a netstring representation of the data in the buffer.
func (n *Netstring) Marshal() ([]byte, error) {
	if !n.complete {
		return nil, Incomplete
	}
	var out bytes.Buffer
	out.WriteString(strconv.Itoa(len(n.buffer)))
	out.WriteRune(':')
	out.Write(n.buffer)
	out.WriteRune(',')
	return out.Bytes(), nil
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

		if length > n.maxLength {
			return TooLarge
		}

		n.buffer = make([]byte, 0, length) // capacity stores the length
	}
	if len(n.buffer) < cap(n.buffer) {
		// slice n.buffer to the part between len and cap
		dest := n.buffer[len(n.buffer):cap(n.buffer)]
		var count int
		count, err = input.Read(dest)

		// slice n.buffer to add on count bytes
		if count > 0 {
			n.buffer = n.buffer[:len(n.buffer)+count]
		}

		switch {
		case err == io.EOF: // we still expect to read a comma, so EOF here is always incomplete
			return Incomplete
		case err != nil:
			return err
		case len(n.buffer) < cap(n.buffer):
			return Incomplete
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
	for {
		digit := make([]byte, 1)
		_, err := input.Read(digit)
		switch {
		case err == io.EOF:
			return -1, Incomplete
		case err != nil:
			return -1, err
		}
		switch rune(digit[0]) {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			n.digits = append(n.digits, digit[0])
			if len(n.digits) > 10 {
				return -1, Garbled
			}
		case ':':
			length, err := strconv.Atoi(string(n.digits))
			if err != nil {
				return -1, Garbled
			}
			return length, nil
		default:
			return -1, Garbled
		}
	}
}

func (n *Netstring) readComma(input io.Reader) error {
	c, _, err := bufio.NewReaderSize(input, 1).ReadRune()
	switch {
	case err == io.EOF:
		return Incomplete
	case err != nil:
		return err
	case c == ',':
		return nil
	default:
		return Garbled
	}
}
