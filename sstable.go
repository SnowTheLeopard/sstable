package sstable

import (
	"bytes"
	"errors"
	"io"
	"sort"
)

var ErrKeyNotFound = errors.New("specified key wasn't found in index")

// SSTable struct returned after reading .sst file
type SSTable struct {
	rws   io.ReadWriteSeeker
	index map[string]uint32
}

// WriteMap converts hashmap to blocks and writes to specified io.Writer
func WriteMap(blocks map[string]string, w io.Writer) error {
	b := NewBlocks()
	for k, v := range blocks {
		b = b.AppendBlock([]byte(k), []byte(v))
	}

	return WriteTable(b, w)
}

// WriteTable write sstable to specified io.Writer
func WriteTable(b Blocks, w io.Writer) error {
	sort.Slice(b, func(i, j int) bool {
		return bytes.Compare(b[i].key, b[j].key) < 0
	})

	var (
		indexOffset int = 8
		blocksBytes []byte
		indexBytes  []byte
	)

	index := NewIndex()

	for _, block := range b {
		b, err := block.marshalBlock()
		if err != nil {
			return err
		}

		index = index.AppendRow(block.key, uint32(indexOffset))

		indexOffset += len(b)
		blocksBytes = append(blocksBytes, b...)
	}

	fileOff := indexOffset
	for _, i := range index {
		bi, err := i.MarshalIndexRow()
		if err != nil {
			return err
		}

		fileOff += len(bi)
		indexBytes = append(indexBytes, bi...)
	}

	h := NewHeader(uint32(fileOff), uint32(indexOffset))

	bh, err := h.MarshalBinary()
	if err != nil {
		return err
	}

	//write header
	if _, err := w.Write(bh); err != nil {
		return err
	}

	//write blocks
	if _, err := w.Write(blocksBytes); err != nil {
		return err
	}

	//write index at the EOF
	if _, err := w.Write(indexBytes); err != nil {
		return err
	}

	return nil
}

// NewTable returns new SSTable struct built from .sst file
// or another source which contains same binary format
// and implements io.ReadWriteSeeker
func NewTable(rws io.ReadWriteSeeker) (*SSTable, error) {
	_, err := rws.Seek(0, 0)
	if err != nil {
		return nil, err
	}

	header, err := readHeader(rws)
	if err != nil {
		return nil, err
	}

	if _, err := rws.Seek(int64(header.indexOffset), 0); err != nil {
		return nil, err
	}

	index, err := readIndex(rws, int(header.fileLen-header.indexOffset))
	if err != nil {
		return nil, err
	}

	sst := &SSTable{
		rws:   rws,
		index: make(map[string]uint32),
	}

	for _, row := range index {
		key := string(row.key)
		sst.index[key] = row.blockpos
	}

	return sst, nil
}

// Search searches for specified key in current sstable
func (sst *SSTable) Search(key string) (string, error) {
	pos, ok := sst.index[key]
	if !ok {
		return "", ErrKeyNotFound
	}

	block, err := readBlock(sst.rws, int64(pos))
	if err != nil {
		return "", err
	}

	return string(block.value), nil
}
