package misc

import (
	"encoding/binary"
	"encoding/json"
	"io"
	"sync/atomic"
)

type Buffer struct {
	idx       int
	len       int
	buffer    []byte
	_refCount int32
	_pool     *BufferPool
}

func (buf *Buffer) Reset() {
	buf.idx = 0
	buf.len = 0
}

func (buf *Buffer) AddRef() *Buffer {
	atomic.AddInt32(&buf._refCount, 1)
	return buf
}

func (buf *Buffer) Recycle() {
	if atomic.AddInt32(&buf._refCount, -1) == 0 {
		buf._pool._pool.Put(buf)
	}
}

func (buf *Buffer) Read(p []byte) (n int, err error) {
	readLen := buf.len - buf.idx
	if readLen > cap(p) {
		readLen = cap(p)
	}
	if readLen > 0 {
		copy(p, buf.buffer[buf.idx:buf.idx+readLen])
		buf.idx += readLen
	}
	return readLen, nil
}

func (buf *Buffer) Write(p []byte) (n int, err error) {
	writeLen := len(p)
	leftSpace := len(buf.buffer) - buf.idx
	if writeLen > leftSpace {
		writeLen = leftSpace
	}
	if writeLen == 0 {
		return 0, nil
	}
	copy(buf.buffer[buf.idx:], p[:writeLen])
	buf.idx += writeLen
	if buf.idx > buf.len {
		buf.len = buf.idx
	}
	return writeLen, nil
}

func (buf *Buffer) Seek(offset int64, whence int) (int64, error) {
	var abs int
	switch whence {
	case 0:
		abs = int(offset)
	case 1:
		abs = buf.idx + int(offset)
	case 2:
		abs = buf.len + int(offset)
	default:
		return 0, nil
	}
	if abs < 0 {
		abs = 0
	}
	if abs > buf.len {
		abs = buf.len
	}
	buf.idx = abs
	return int64(abs), nil
}

func (buf *Buffer) WrtieJson(obj interface{}) (n int, err error) {
	oriPos := buf.idx
	encoder := json.NewEncoder(buf)
	if e := encoder.Encode(obj); e != nil {
		return 0, e
	}
	return buf.idx - oriPos, nil
}

func (buf *Buffer) ReadJson(obj interface{}) (n int, err error) {
	oriPos := buf.idx
	decoder := json.NewDecoder(buf)
	if e := decoder.Decode(obj); e != nil {
		return 0, e
	}
	return buf.idx - oriPos, nil
}

func (buf *Buffer) Len() int {
	return buf.len - buf.idx
}

func (buf *Buffer) Cap() int {
	return len(buf.buffer) - buf.idx
}

func (buf *Buffer) Pos() int {
	return buf.idx
}

func (buf *Buffer) Bytes() []byte {
	return buf.buffer[0:buf.len]
}

func (buf *Buffer) ReadByte() (byte, error) {
	if buf.idx >= buf.len {
		return 0, io.EOF
	}
	b := buf.buffer[buf.idx]
	buf.idx++
	return b, nil
}

func (buf *Buffer) WriteByte(b byte) error {
	if buf.idx >= len(buf.buffer) {
		return io.ErrShortBuffer
	}
	buf.buffer[buf.idx] = b
	buf.idx++
	if buf.idx > buf.len {
		buf.len = buf.idx
	}
	return nil
}

func (buf *Buffer) ReadLittleEndian(data any) error {
	if buf.idx+4 > buf.len {
		return io.EOF
	}
	err := binary.Read(buf, binary.LittleEndian, data)
	if err != nil {
		return err
	}
	buf.idx += 4
	return nil
}

func (buf *Buffer) ReadBigEndian(data any) error {
	if buf.idx+4 > buf.len {
		return io.EOF
	}
	err := binary.Read(buf, binary.BigEndian, data)
	if err != nil {
		return err
	}
	buf.idx += 4
	return nil
}

func (buf *Buffer) WriteLittleEndian(v any) error {
	if buf.idx+4 > len(buf.buffer) {
		return io.ErrShortBuffer
	}
	err := binary.Write(buf, binary.LittleEndian, v)
	if err != nil {
		return err
	}
	buf.idx += 4
	if buf.idx > buf.len {
		buf.len = buf.idx
	}
	return nil
}

func (buf *Buffer) WriteBigEndian(v any) error {
	if buf.idx+4 > len(buf.buffer) {
		return io.ErrShortBuffer
	}
	err := binary.Write(buf, binary.BigEndian, v)
	if err != nil {
		return err
	}
	buf.idx += 4
	if buf.idx > buf.len {
		buf.len = buf.idx
	}
	return nil
}
