package amino

import (
	"bytes"
	"sync"
)

type BufferPool struct {
	pool *sync.Pool
}

func NewBufferPool() *BufferPool {
	return &BufferPool{
		pool: &sync.Pool{
			New: func() interface{} {
				return &bytes.Buffer{}
			},
		},
	}
}

func (p *BufferPool) Get() *bytes.Buffer {
	ret := p.pool.Get().(*bytes.Buffer)
	ret.Reset()
	return ret
}

func (p *BufferPool) Put(buf *bytes.Buffer) {
	p.pool.Put(buf)
}
