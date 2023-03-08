package bytesrw

import (
	"fmt"
	"io"
)

// BytesReadWriteSeeker implements Read, Write,
// Seek and ReadByte functions for byte slice.
type BytesReadWriteSeeker struct {
	data []byte
	pos  int
}

func NewBytesReadWriteSeeker(data []byte) *BytesReadWriteSeeker {
	return &BytesReadWriteSeeker{data: data, pos: 0}
}

// Seek moves the pointer to the specified position.
func (brws *BytesReadWriteSeeker) Seek(offset int64, whence int) (int64, error) {
	var start int

	switch whence {
	case io.SeekStart:
		start = 0

	case io.SeekCurrent:
		start = brws.pos

	case io.SeekEnd:
		start = len(brws.data)

	default:
		return -1, fmt.Errorf("option not defined")
	}

	newPos := start + int(offset)

	switch {
	case newPos < 0:
		newPos = 0

	case newPos > len(brws.data):
		newPos = len(brws.data)
	}

	brws.pos = newPos

	return int64(brws.pos), nil
}

// Write writes the data to the inner buffer.
func (brws *BytesReadWriteSeeker) Write(p []byte) (n int, err error) {
	offset := len(p)
	brws.data = append(brws.data, p...)
	brws.pos += offset

	return offset, nil
}

// ReadByte reads exactly 1 byte from the inner buffer.
func (brws *BytesReadWriteSeeker) ReadByte() (byte, error) {
	if brws.pos == len(brws.data) {
		return 0, io.EOF
	}

	value := brws.data[brws.pos]
	brws.pos++

	return value, nil
}

// Read reads the bytes from the inner buffer to the provided buffer.
func (brws *BytesReadWriteSeeker) Read(p []byte) (n int, err error) {
	if brws.pos == len(brws.data) {
		return -1, io.EOF
	}

	offset := len(p)

	if offset > len(brws.data)-brws.pos {
		offset = len(brws.data) - brws.pos
	}

	copy(p, brws.data[brws.pos:brws.pos+offset])
	brws.pos += offset

	return offset, nil
}
