package sstable

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

const headerLen int = 8

type Header struct {
	fileLen     uint32
	indexOffset uint32
}

func NewHeader(fl, io uint32) Header {
	return Header{
		fileLen:     fl,
		indexOffset: io,
	}
}

func (h *Header) MarshalBinary() ([]byte, error) {
	buf := bytes.NewBuffer([]byte{})

	if err := binary.Write(buf, binary.BigEndian, h); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func UnmarshalHeader(data []byte) (Header, error) {
	h := Header{}

	if len(data) < 8 {
		return h, fmt.Errorf("header data is less then required")
	}

	flen := binary.BigEndian.Uint32(data[:4])
	ioff := binary.BigEndian.Uint32(data[4:8])

	h.fileLen = flen
	h.indexOffset = ioff

	return h, nil
}

func readHeader(r io.Reader) (Header, error) {
	h := Header{}

	bhead := make([]byte, headerLen)
	if _, err := r.Read(bhead); err != nil {
		return h, err
	}

	h, err := UnmarshalHeader(bhead)
	if err != nil {
		return h, err
	}

	return h, nil
}
