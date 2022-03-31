package sstable

import (
	"bytes"
	"encoding/binary"
	"io"
)

type Blocks []BlockRow

func NewBlocks() Blocks {
	return []BlockRow{}
}

func (b Blocks) AppendBlock(k, v []byte) Blocks {
	block := BlockRow{
		key:   k,
		value: v,
	}

	return append(b, block)
}

type BlockRow struct {
	key   []byte
	value []byte
}

func (br *BlockRow) marshalBlock() ([]byte, error) {
	buf := bytes.NewBuffer([]byte{})

	if err := binary.Write(buf, binary.BigEndian, uint32(len(br.key))); err != nil {
		return nil, err
	}

	if err := binary.Write(buf, binary.BigEndian, uint32(len(br.value))); err != nil {
		return nil, err
	}

	if _, err := buf.Write(br.key); err != nil {
		return nil, err
	}

	if _, err := buf.Write(br.value); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func readBlock(rs io.ReadSeeker, blockPos int64) (BlockRow, error) {
	br := BlockRow{}

	_, err := rs.Seek(blockPos, 0)
	if err != nil {
		return br, err
	}

	meta := make([]byte, 8)
	if _, err := rs.Read(meta); err != nil {
		return br, err
	}

	kl := binary.BigEndian.Uint32(meta[:4])
	vl := binary.BigEndian.Uint32(meta[4:8])

	_, err = rs.Seek(blockPos+8, 0)
	if err != nil {
		return br, err
	}

	data := make([]byte, kl+vl)
	if _, err := rs.Read(data); err != nil {
		return br, err
	}

	br.key = make([]byte, kl)
	copy(br.key, data[:kl])

	br.value = make([]byte, vl)
	copy(br.value, data[kl:])

	return br, nil
}
