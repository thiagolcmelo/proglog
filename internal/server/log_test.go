package server

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestLogAppend(t *testing.T) {
	s := NewLog()
	expected := Record{Value: []byte("oyh5re2w"), Offset: 0}
	record := Record{Value: []byte("oyh5re2w")}

	offset, err := s.Append(record)
	if err != nil {
		t.Errorf("err is not nil: %v", err)
	}
	if offset != 0 {
		t.Errorf("offset is not 0: %d", offset)
	}

	actual := s.records[0]
	if diff := cmp.Diff(actual, expected); diff != "" {
		t.Error(diff)
	}
}

func TestLogRead(t *testing.T) {
	s := NewLog()
	expected := Record{Value: []byte("oyh5re2w"), Offset: 0}
	record := Record{Value: []byte("oyh5re2w")}
	_, err := s.Append(record)
	if err != nil {
		t.Errorf("err is not nil: %v", err)
	}

	actual, err := s.Read(0)
	if err != nil {
		t.Errorf("err is not nil: %v", err)
	}
	if diff := cmp.Diff(actual, expected); diff != "" {
		t.Error(diff)
	}
}

func TestReadInvalidOffsetFails(t *testing.T) {
	s := NewLog()
	record := Record{Value: []byte("oyh5re2w")}
	_, err := s.Append(record)
	if err != nil {
		t.Errorf("err is not nil: %v", err)
	}

	_, err = s.Read(1)
	if err == nil {
		t.Error("expected error to be not nil")
	}
}
