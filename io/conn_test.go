package io

import (
	"io"
	"sync"
	"testing"
)

const (
	bufferLength   = 8 * 1024           // buffer 长度
	transferLength = 1024 * 1024 * 1024 // GB, 每次测试传输的数据长度
)

func BenchmarkMemoryConnType(b *testing.B) {

	testReader := func(wg *sync.WaitGroup, conn io.ReadWriteCloser) {
		defer wg.Done()

		buffer := make([]byte, bufferLength)
		hasRead := 0
		for {
			if hasRead >= transferLength {
				break
			}

			n, err := conn.Read(buffer[:bufferLength])
			if err != nil {
				b.Errorf("read error: %s", err)
			}
			hasRead += n
		}
	}

	testWriter := func(wg *sync.WaitGroup, conn io.ReadWriteCloser) {
		defer wg.Done()

		buffer := make([]byte, bufferLength)
		hasWrite := 0
		for {
			if hasWrite >= transferLength {
				break
			}

			n, err := conn.Write(buffer[:bufferLength])
			if err != nil {
				b.Errorf("read error: %s", err)
			}
			hasWrite += n
		}
	}

	for i := 0; i < b.N; i++ {
		client, server := NewMemoryConnPeer()
		wg := &sync.WaitGroup{}
		wg.Add(4)
		go testReader(wg, client)
		go testWriter(wg, server)
		go testReader(wg, server)
		go testWriter(wg, client)
		wg.Wait()
	}
}
