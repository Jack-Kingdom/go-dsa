package buffer

import (
	"math/rand"
	"testing"
	"time"
)

func init () {
	rand.Seed(time.Now().UnixNano())
}

const (
	maxSize = DefaultAllocatorMaxSize
)

func TestGet(t *testing.T) {
	size := rand.Intn(maxSize)

	buffer := Get(size)
	defer Put(buffer)

	t.Logf("buffer: reqeust len %d, cap %d", size, cap(buffer))
	if cap(buffer) < size {
		t.Errorf("buffer cap less than request len")
	}
}