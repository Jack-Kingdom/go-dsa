package buffer

import (
	"errors"
	"sync"
)

var (
	// ErrCircularBufferEmpty is returned when the buffer is empty.
	ErrCircularBufferEmpty = errors.New("circular buffer is empty")
	// ErrCircularBufferFull is returned when the buffer is full.
	ErrCircularBufferFull = errors.New("circular buffer is full")
	// ErrTooManyDataToWrite is returned when the buffer available is less than ready to write.
	ErrTooManyDataToWrite = errors.New("circular buffer too many data to write")
)

// CircularBuffer is a thread-safe circular buffer.
type CircularBuffer struct {
	size   int
	buf    []byte
	r, w   int // next position to read & write
	length int // current used length
	mu     sync.Mutex
}

func NewCircularBuffer(buf []byte) *CircularBuffer {
	return &CircularBuffer{
		size: len(buf),
		buf:  buf,
	}
}

func (circular *CircularBuffer) Len() int {
	return circular.length
}

func (circular *CircularBuffer) Read(data []byte) (n int, err error) {
	circular.mu.Lock()
	defer circular.mu.Unlock()

	if circular.length == 0 {
		return 0, ErrCircularBufferEmpty
	}

	if circular.r < circular.w {
		n = copy(data, circular.buf[circular.r:circular.w])
		circular.r = circular.r + n
	} else {
		n = copy(data, circular.buf[circular.r:])
		n += copy(data[n:], circular.buf[:circular.w])
		circular.r = circular.r + n - circular.size
	}

	circular.length -= n
	return n, nil
}

func (circular *CircularBuffer) Write(data []byte) (n int, err error) {
	if len(data) == 0 {
		return 0, nil
	}

	circular.mu.Lock()
	defer circular.mu.Unlock()

	if circular.length >= circular.size {
		return 0, ErrCircularBufferFull
	}

	avail := circular.size - circular.length
	if avail < len(data) {
		return 0, ErrTooManyDataToWrite
	}

	if circular.r < circular.w {
		if circular.size-circular.w >= len(data) {
			n = copy(circular.buf[circular.w:], data)
			circular.w = circular.w + n
		} else {
			n = copy(circular.buf[circular.w:], data)
			n += copy(circular.buf[:circular.r], data[n:])
			circular.w = circular.w + n - circular.size
		}
	} else {
		n = copy(circular.buf[circular.w:], data)
		circular.w = circular.w + n
	}

	circular.length += n
	return n, nil
}