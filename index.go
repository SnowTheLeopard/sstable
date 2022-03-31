package sstable

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

type Index []IndexRow

func NewIndex() Index {
	return []IndexRow{}
}

func (i Index) AppendRow(key []byte, pos uint32) Index {
	row := IndexRow{
		key:      key,
		blockpos: pos,
	}

	return append(i, row)
}

type IndexRow struct {
	key      []byte
	blockpos uint32
}

func (ir *IndexRow) MarshalIndexRow() ([]byte, error) {
	buf := bytes.NewBuffer([]byte{})

	if err := binary.Write(buf, binary.BigEndian, ir.blockpos); err != nil {
		return nil, err
	}

	if err := binary.Write(buf, binary.BigEndian, uint32(len(ir.key))); err != nil {
		return nil, err
	}

	if _, err := buf.Write(ir.key); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func UnmarshalIndexRow(data []byte) (IndexRow, int, error) {
	ir := IndexRow{}

	if len(data) < 4 {
		return ir, 0, fmt.Errorf("index row is less then required")
	}

	ir.blockpos = binary.BigEndian.Uint32(data[:4])

	kl := binary.BigEndian.Uint32(data[4:8])
	ir.key = make([]byte, kl)
	copy(ir.key, data[8:8+kl])

	return ir, int(8 + kl), nil
}

func readIndex(r io.Reader, size int) (Index, error) {
	index := make([]byte, size)
	if _, err := r.Read(index); err != nil {
		return nil, err
	}

	var (
		read = 0
		rows []IndexRow
	)

	for read < len(index) {
		ir, br, err := UnmarshalIndexRow(index[read:])
		if err != nil {
			return nil, err
		}

		read += br
		rows = append(rows, ir)
	}

	return rows, nil
}
