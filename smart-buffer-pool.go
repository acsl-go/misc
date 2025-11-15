package misc

import "sync"

type SmartBufferPool struct {
	bufferSizes []uint
	_pool       []*sync.Pool
}

var (
	defaultSmartBufferPoolSizes = []uint{
		128 * 1024,
		512 * 1024,
		1 * 1024 * 1024,
		8 * 1024 * 1024,
		16 * 1024 * 1024,
		32 * 1024 * 1024,
		64 * 1024 * 1024,
		128 * 1024 * 1024,
	}
)

func NewSmartBufferPool() *SmartBufferPool {
	return NewSmartBufferPoolEx(defaultSmartBufferPoolSizes)
}

func NewSmartBufferPoolEx(sizes []uint) *SmartBufferPool {
	// 128k, 512k, 1M, 8M, 16M, 32M, 64M, 128M
	if sizes == nil {
		sizes = defaultSmartBufferPoolSizes
	}
	pool := &SmartBufferPool{
		bufferSizes: make([]uint, len(sizes)),
		_pool:       make([]*sync.Pool, len(sizes)),
	}
	for i, sz := range sizes {
		pool.bufferSizes[i] = sz
		pool._pool[i] = &sync.Pool{
			New: func() interface{} {
				return NewBuffer(sz)
			},
		}
	}
	return pool
}

func (p *SmartBufferPool) Get(size uint) *Buffer {
	for i, sz := range p.bufferSizes {
		if size <= sz {
			buf := p._pool[i].Get().(*Buffer)
			buf._refCount = 1
			buf._pool = p._pool[i]
			buf.Resize(size, false)
			buf.Reset()
			return buf
		}
	}
	buf := NewBuffer(size)
	buf._refCount = 1
	buf._pool = nil
	return buf
}
