package buffer

import (
	"log"
	"sync"
)

const (
	DefaultAllocatorMaxSize = 128 * 1024 // 128KB
)

var (
	defaultAllocator     *Allocator
	defaultAllocatorOnce sync.Once
)

// Allocator is a memory allocator.
// 每次分配的内存都满足 2 的幂，计算最为接近的内存大小，保证浪费的空间不超过 50%
type Allocator struct {
	maxSize int
	shiftLimit int
	pools   []sync.Pool
}

func NewAllocator(maxSize int) *Allocator {
	allocator := new(Allocator)
	allocator.maxSize = maxSize
	allocator.shiftLimit = calShiftLimit(maxSize)

	allocator.pools = make([]sync.Pool, allocator.shiftLimit)
	for i := range allocator.pools {
		j := i	// 避免对 i 的引用
		allocator.pools[i].New = func() interface{} {
			return make([]byte, 1<<uint(j))
		}
	}

	return allocator
}

func calShiftLimit(size int) int {
	if size <= 0 {
		return 0
	}

	shift := 0
	for size > 0 {
		size >>= 1
		shift++
	}

	return shift
}

func (allocator *Allocator) Get(size int) []byte {
	if size > allocator.maxSize {
		return nil
	}

	for i := 0; i < allocator.shiftLimit; i++ {
		if (1 << i) > size {
			return allocator.pools[i].Get().([]byte)[:]
		}
	}
	log.Println("allocator.Get: request size larger than shiftLimit")
	return nil
}

func (allocator *Allocator) Put(buf []byte) {
	if cap(buf) > allocator.maxSize {
		return
	}

	for i := 0; i < allocator.shiftLimit; i++ {
		if (1 << i) == cap(buf) {
			allocator.pools[i].Put(buf)
			return
		}
	}
	log.Println("WARN: buffer.Put: buf size miss match")
}

func Get(size int) []byte {
	defaultAllocatorOnce.Do(func() {
		defaultAllocator = NewAllocator(DefaultAllocatorMaxSize)
	})
	return defaultAllocator.Get(size)
}

func Put(buf []byte) {
	defaultAllocatorOnce.Do(func() {
		defaultAllocator = NewAllocator(DefaultAllocatorMaxSize)
	})
	defaultAllocator.Put(buf)
}
