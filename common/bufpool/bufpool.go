// Package bufpool provides a sync.Pool of bytes.Buffer.
package bufpool

import (
	"bytes"
	"sync"
)

const maxSize = 1 << 16 // 64 KiB. See https://github.com/golang/go/issues/23199

var p = &sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

// Get gets a buffer from the pool, resets it and returns it.
func Get() *bytes.Buffer {
	buf := p.Get().(*bytes.Buffer)
	buf.Reset()
	return buf
}

// Put puts the buffer to the pool.
// WARNING: the call MUST NOT reuse the buffer's content after this call.
func Put(buf *bytes.Buffer) {
	if buf.Cap() <= maxSize {
		p.Put(buf)
	}
}
