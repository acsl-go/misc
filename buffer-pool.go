package misc

import "sync"

type BufferPool struct {
	MaxBufferSize uint
	_pool         sync.Pool
}

func NewBufferPool(maxBufferSize uint) *BufferPool {
	return &BufferPool{
		MaxBufferSize: maxBufferSize,
		_pool: sync.Pool{
			New: func() interface{} {
				return &Buffer{
					buffer: make([]byte, maxBufferSize),
				}
			},
		},
	}
}

func (p *BufferPool) Get() *Buffer {
	buf := p._pool.Get().(*Buffer)
	buf._refCount = 1
	buf._pool = p
	buf.Reset()
	return buf
}
