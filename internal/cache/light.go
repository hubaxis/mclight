package cache

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

var (
	crlf          = []byte("\r\n")
	resultStored  = []byte("STORED\r\n")
	resultDeleted = []byte("DELETED\r\n")
	resultOk      = []byte("OK\r\n")
	resultEnd     = []byte("END\r\n")
)

type light struct {
	conn   net.Conn
	connMu sync.Mutex
}

// New instance of light
func New(endpoint string) (Cache, error) {
	conn, err := net.Dial("tcp", endpoint)
	if err != nil {
		return nil, err
	}
	return &light{
		conn:   conn,
		connMu: sync.Mutex{},
	}, nil
}

// Get value by key
func (l *light) Get(key string) ([]byte, error) {
	l.connMu.Lock()
	defer l.connMu.Unlock()
	buffer := bufio.NewReadWriter(bufio.NewReader(l.conn), bufio.NewWriter(l.conn))
	_, err := fmt.Fprintf(buffer, "gets %s\r\n", key)
	if err != nil {
		return nil, err
	}
	err = buffer.Flush()
	if err != nil {
		return nil, err
	}
	line, err := buffer.ReadSlice('\n')
	if err != nil {
		return nil, err
	}
	if bytes.Equal(line, resultEnd) {
		return nil, nil
	}
	size, err := scan(line)
	if err != nil {
		return nil, err
	}
	value := make([]byte, size+2)
	_, err = io.ReadFull(buffer, value)
	if err != nil {
		return nil, err
	}
	if !bytes.HasSuffix(value, crlf) {
		return nil, fmt.Errorf("unexpected respose %v", value)
	}
	return value[:size], nil
}

// Set value with key and expiration
func (l *light) Set(key string, value []byte, expiration time.Duration) error {
	l.connMu.Lock()
	defer l.connMu.Unlock()
	buffer := bufio.NewReadWriter(bufio.NewReader(l.conn), bufio.NewWriter(l.conn))
	_, err := fmt.Fprintf(buffer, "%s %s %d %d %d\r\n",
		"set", key, 0, int(expiration.Seconds()), len(value))
	if err != nil {
		return err
	}
	_, err = buffer.Write(value)
	if err != nil {
		return err
	}
	_, err = buffer.Write(crlf)
	if err != nil {
		return err
	}
	err = buffer.Flush()
	if err != nil {
		return err
	}
	line, err := buffer.ReadSlice('\n')
	if err != nil {
		return err
	}
	if bytes.Equal(line, resultStored) {
		return nil
	}
	return fmt.Errorf("unexpected response %s", string(line))
}

// Delete value by key
func (l *light) Delete(key string) error {
	l.connMu.Lock()
	defer l.connMu.Unlock()
	buffer := bufio.NewReadWriter(bufio.NewReader(l.conn), bufio.NewWriter(l.conn))
	_, err := fmt.Fprintf(buffer, "delete %s\r\n", key)
	if err != nil {
		return err
	}
	if err := buffer.Flush(); err != nil {
		return err
	}
	line, err := buffer.ReadSlice('\n')
	if bytes.Equal(line, resultOk) {
		return nil
	}
	if bytes.Equal(line, resultDeleted) {
		return nil
	}
	return fmt.Errorf("can't delete item")
}

// Close the client
func (l *light) Close() error {
	l.connMu.Lock()
	defer l.connMu.Unlock()
	return l.conn.Close()
}

func scan(line []byte) (int, error) {
	pattern := "VALUE %s %d %d %d\r\n"
	var key string
	var flags uint32
	var size int
	var casID uint64
	dest := []interface{}{&key, &flags, &size, &casID}
	n, err := fmt.Sscanf(string(line), pattern, dest...)
	if err != nil || n != len(dest) {
		return -1, fmt.Errorf("unexpected response %v", line)
	}
	return size, nil
}
