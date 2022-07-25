package io

import (
	stdio "io"
)

// MemoryConnType 是一个内存连接类型，通常用于测试
type MemoryConnType struct {
	reader *stdio.PipeReader
	writer *stdio.PipeWriter
}

func (conn *MemoryConnType) Read(b []byte) (int, error) {
	return conn.reader.Read(b)
}

func (conn *MemoryConnType) Write(b []byte) (int, error) {
	return conn.writer.Write(b)
}

func (conn *MemoryConnType) Close() error {
	return conn.writer.Close()
}

func NewMemoryConnPeer() (*MemoryConnType, *MemoryConnType) {
	clientReader, serverWriter := stdio.Pipe()
	serverReader, clientWriter := stdio.Pipe()
	return &MemoryConnType{clientReader, clientWriter}, &MemoryConnType{serverReader, serverWriter}
}
