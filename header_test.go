package sstable_test

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/SnowTheLeopard/sstable"
)

func Test_MarshalHeader(t *testing.T) {
	test := struct {
		name string
		arg  sstable.Header
		want []byte
	}{
		name: "test header marshaling",
		arg:  sstable.NewHeader(2, 1),
		want: []byte{0, 0, 0, 2, 0, 0, 0, 1},
	}

	t.Run(test.name, func(t *testing.T) {
		data, err := test.arg.MarshalBinary()
		if err != nil {
			t.Errorf("failed to marshal header: %w", err)
		}

		if eq := bytes.Compare(test.want, data); eq != 0 {
			t.Error("received another binary data")
		}
	})
}

func Test_UnmarshalHeader(t *testing.T) {
	test := struct {
		name string
		arg  []byte
		want sstable.Header
	}{
		name: "test header unmarshal",
		arg:  []byte{0, 0, 0, 2, 0, 0, 0, 1},
		want: sstable.NewHeader(2, 1),
	}

	t.Run(test.name, func(t *testing.T) {
		h, err := sstable.UnmarshalHeader(test.arg)
		if err != nil {
			t.Errorf("failed to unmarshal header: %w", err)
		}

		if eq := reflect.DeepEqual(test.want, h); !eq {
			t.Error("header not equal")
		}
	})
}
