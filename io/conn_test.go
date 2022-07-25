package io

import (
	"context"
	"io"
	"sync"
	"testing"
	"time"
)

const (
	bufferLength = 8 * 1024 // buffer 长度
	Mb           = 1024 * 1024
)

var (
	mutex sync.Mutex
)

func BenchmarkMemoryConnType(b *testing.B) {
	testReader := func(ctx context.Context, conn io.ReadWriteCloser) {
		buffer := make([]byte, bufferLength)
		hasRead := 0
		for {
			select {
			case <-ctx.Done():
				mutex.Lock()
				b.ReportMetric(float64(hasRead)/Mb, "MB")
				mutex.Unlock()
				return
			default:
				n, err := conn.Read(buffer[:bufferLength])
				if err != nil {
					b.Errorf("read error: %s", err)
				}
				hasRead += n
			}
		}
	}

	testWriter := func(ctx context.Context, conn io.ReadWriteCloser) {
		buffer := make([]byte, bufferLength)
		hasWrite := 0
		for {
			select {
			case <-ctx.Done():
				return
			default:
				n, err := conn.Write(buffer[:bufferLength])
				if err != nil {
					b.Errorf("read error: %s", err)
				}
				hasWrite += n
			}
		}
	}

	for i := 0; i < b.N; i++ {
		b.ReportAllocs()

		client, server := NewMemoryConnPeer()
		ctx, cancel := context.WithTimeout(context.TODO(), time.Second)
		go testReader(ctx, client)
		go testWriter(ctx, server)
		go testReader(ctx, server)
		go testWriter(ctx, client)

		time.Sleep(1 * time.Second)
		cancel()
	}
}
