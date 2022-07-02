package server

import (
	"fmt"
	"sync"
)

// Log is an in memory log
type Log struct {
	mu      sync.Mutex
	records []Record
}

// NewLog produces a new Log
func NewLog() *Log {
	return &Log{}
}

// Append adds a new entry to the in memory log
func (c *Log) Append(record Record) (uint64, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	record.Offset = uint64(len(c.records))
	c.records = append(c.records, record)
	return record.Offset, nil
}

// Read return a record from the in memory log by its offset
func (c *Log) Read(offset uint64) (Record, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if offset > uint64(len(c.records)-1) {
		return Record{}, ErrOffsetNotFound
	}
	return c.records[offset], nil
}

// ErrOffsetNotFound informs that a record with such offset cannot be found
var ErrOffsetNotFound = fmt.Errorf("offset not found")

// Record represents a log entry
type Record struct {
	Value  []byte `json:"value"`
	Offset uint64 `json:"offset"`
}
