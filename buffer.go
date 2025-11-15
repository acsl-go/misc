package misc

import (
	"encoding/binary"
	"encoding/json"
	"io"
	"sync"
	"sync/atomic"
)

const (
	SEEK_SET = 0
	SEEK_CUR = 1
	SEEK_END = 2
)

type Buffer struct {
	idx       int
	len       int
	buffer    []byte
	Tag       int
	_refCount int32
	_pool     *sync.Pool
}

func NewBuffer(sz uint) *Buffer {
	return &Buffer{
		idx:       0,
		len:       0,
		buffer:    make([]byte, sz),
		_refCount: 1,
		_pool:     nil,
	}
}

func (buf *Buffer) Reset() {
	buf.idx = 0
	buf.len = 0
}

func (buf *Buffer) AddRef() *Buffer {
	atomic.AddInt32(&buf._refCount, 1)
	return buf
}

func (buf *Buffer) Resize(sz uint, keepData bool) {
	if sz > uint(cap(buf.buffer)) {
		newBuffer := make([]byte, sz)
		if keepData {
			copy(newBuffer, buf.buffer[:buf.len])
		}
		buf.buffer = newBuffer
	}
}

func (buf *Buffer) Release() {
	if atomic.AddInt32(&buf._refCount, -1) == 0 {
		if buf._pool != nil {
			buf._pool.Put(buf)
		}
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

func (buf *Buffer) WriteJson(obj interface{}) (n int, err error) {
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

func (buf *Buffer) Data() []byte {
	if buf.idx >= buf.len {
		return nil
	}
	return buf.buffer[buf.idx:buf.len]
}

func (buf *Buffer) Bytes() []byte {
	return buf.buffer[0:buf.len]
}

func (buf *Buffer) Buffer() []byte {
	return buf.buffer
}

func (buf *Buffer) SetDataLen(len int) {
	buf.idx = 0
	buf.len = len
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

func (buf *Buffer) ReadString() (string, error) {
	if buf.idx+4 > buf.len {
		return "", io.EOF
	}

	n, _ := buf.ReadUint32BE()
	if buf.idx+int(n) > buf.len {
		return "", io.EOF
	}

	s := string(buf.buffer[buf.idx : buf.idx+int(n)])
	buf.idx += int(n)
	return s, nil
}

func (buf *Buffer) MustReadString() string {
	s, err := buf.ReadString()
	if err != nil {
		panic(err)
	}
	return s
}

func (buf *Buffer) WriteString(s string) error {
	if e := buf.WriteUint32BE(uint32(len(s))); e != nil {
		return e
	}
	if _, e := buf.Write([]byte(s)); e != nil {
		return e
	}
	return nil
}

func (buf *Buffer) ReadCString() (string, error) {
	start := buf.idx
	for {
		if start >= buf.len {
			return "", io.EOF
		}
		if buf.buffer[start] == 0 {
			s := string(buf.buffer[buf.idx:start])
			buf.idx = start + 1
			return s, nil
		}
	}
}

func (buf *Buffer) MustReadCString() string {
	s, err := buf.ReadCString()
	if err != nil {
		panic(err)
	}
	return s
}

func (buf *Buffer) WriteCString(s string) error {
	if _, e := buf.Write([]byte(s)); e != nil {
		return e
	}
	if e := buf.WriteByte(0); e != nil {
		return e
	}
	return nil
}

func (buf *Buffer) ReadLittleEndian(data any) error {
	return binary.Read(buf, binary.LittleEndian, data)
}

func (buf *Buffer) ReadBigEndian(data any) error {
	return binary.Read(buf, binary.BigEndian, data)
}

func (buf *Buffer) WriteLittleEndian(v any) error {
	return binary.Write(buf, binary.LittleEndian, v)
}

func (buf *Buffer) WriteBigEndian(v any) error {
	return binary.Write(buf, binary.BigEndian, v)
}

func (buf *Buffer) ReadInt8LE() (int8, error) {
	var v int8
	if err := buf.ReadLittleEndian(&v); err != nil {
		return 0, err
	}
	return v, nil
}

func (buf *Buffer) WriteInt8LE(v int8) error {
	return buf.WriteLittleEndian(v)
}

func (buf *Buffer) ReadUint8LE() (uint8, error) {
	var v uint8
	if err := buf.ReadLittleEndian(&v); err != nil {
		return 0, err
	}
	return v, nil
}

func (buf *Buffer) WriteUint8LE(v uint8) error {
	return buf.WriteLittleEndian(v)
}

func (buf *Buffer) ReadInt16LE() (int16, error) {
	var v int16
	if err := buf.ReadLittleEndian(&v); err != nil {
		return 0, err
	}
	return v, nil
}

func (buf *Buffer) WriteInt16LE(v int16) error {
	return buf.WriteLittleEndian(v)
}

func (buf *Buffer) ReadUint16LE() (uint16, error) {
	var v uint16
	if err := buf.ReadLittleEndian(&v); err != nil {
		return 0, err
	}
	return v, nil
}

func (buf *Buffer) WriteUint16LE(v uint16) error {
	return buf.WriteLittleEndian(v)
}

func (buf *Buffer) ReadInt32LE() (int32, error) {
	var v int32
	if err := buf.ReadLittleEndian(&v); err != nil {
		return 0, err
	}
	return v, nil
}

func (buf *Buffer) WriteInt32LE(v int32) error {
	return buf.WriteLittleEndian(v)
}

func (buf *Buffer) ReadUint32LE() (uint32, error) {
	var v uint32
	if err := buf.ReadLittleEndian(&v); err != nil {
		return 0, err
	}
	return v, nil
}

func (buf *Buffer) WriteUint32LE(v uint32) error {
	return buf.WriteLittleEndian(v)
}

func (buf *Buffer) ReadInt64LE() (int64, error) {
	var v int64
	if err := buf.ReadLittleEndian(&v); err != nil {
		return 0, err
	}
	return v, nil
}

func (buf *Buffer) WriteInt64LE(v int64) error {
	return buf.WriteLittleEndian(v)
}

func (buf *Buffer) ReadUint64LE() (uint64, error) {
	var v uint64
	if err := buf.ReadLittleEndian(&v); err != nil {
		return 0, err
	}
	return v, nil
}

func (buf *Buffer) WriteUint64LE(v uint64) error {
	return buf.WriteLittleEndian(v)
}

func (buf *Buffer) ReadFloat32LE() (float32, error) {
	var v float32
	if err := buf.ReadLittleEndian(&v); err != nil {
		return 0, err
	}
	return v, nil
}

func (buf *Buffer) WriteFloat32LE(v float32) error {
	return buf.WriteLittleEndian(v)
}

func (buf *Buffer) ReadFloat64LE() (float64, error) {
	var v float64
	if err := buf.ReadLittleEndian(&v); err != nil {
		return 0, err
	}
	return v, nil
}

func (buf *Buffer) WriteFloat64LE(v float64) error {
	return buf.WriteLittleEndian(v)
}

func (buf *Buffer) ReadInt8BE() (int8, error) {
	var v int8
	if err := buf.ReadBigEndian(&v); err != nil {
		return 0, err
	}
	return v, nil
}

func (buf *Buffer) WriteInt8BE(v int8) error {
	return buf.WriteBigEndian(v)
}

func (buf *Buffer) ReadUint8BE() (uint8, error) {
	var v uint8
	if err := buf.ReadBigEndian(&v); err != nil {
		return 0, err
	}
	return v, nil
}

func (buf *Buffer) WriteUint8BE(v uint8) error {
	return buf.WriteBigEndian(v)
}

func (buf *Buffer) ReadInt16BE() (int16, error) {
	var v int16
	if err := buf.ReadBigEndian(&v); err != nil {
		return 0, err
	}
	return v, nil
}

func (buf *Buffer) WriteInt16BE(v int16) error {
	return buf.WriteBigEndian(v)
}

func (buf *Buffer) ReadUint16BE() (uint16, error) {
	var v uint16
	if err := buf.ReadBigEndian(&v); err != nil {
		return 0, err
	}
	return v, nil
}

func (buf *Buffer) WriteUint16BE(v uint16) error {
	return buf.WriteBigEndian(v)
}

func (buf *Buffer) ReadInt32BE() (int32, error) {
	var v int32
	if err := buf.ReadBigEndian(&v); err != nil {
		return 0, err
	}
	return v, nil
}

func (buf *Buffer) WriteInt32BE(v int32) error {
	return buf.WriteBigEndian(v)
}

func (buf *Buffer) ReadUint32BE() (uint32, error) {
	var v uint32
	if err := buf.ReadBigEndian(&v); err != nil {
		return 0, err
	}
	return v, nil
}

func (buf *Buffer) WriteUint32BE(v uint32) error {
	return buf.WriteBigEndian(v)
}

func (buf *Buffer) ReadInt64BE() (int64, error) {
	var v int64
	if err := buf.ReadBigEndian(&v); err != nil {
		return 0, err
	}
	return v, nil
}

func (buf *Buffer) WriteInt64BE(v int64) error {
	return buf.WriteBigEndian(v)
}

func (buf *Buffer) ReadUint64BE() (uint64, error) {
	var v uint64
	if err := buf.ReadBigEndian(&v); err != nil {
		return 0, err
	}
	return v, nil
}

func (buf *Buffer) WriteUint64BE(v uint64) error {
	return buf.WriteBigEndian(v)
}

func (buf *Buffer) ReadFloat32BE() (float32, error) {
	var v float32
	if err := buf.ReadBigEndian(&v); err != nil {
		return 0, err
	}
	return v, nil
}

func (buf *Buffer) WriteFloat32BE(v float32) error {
	return buf.WriteBigEndian(v)
}

func (buf *Buffer) ReadFloat64BE() (float64, error) {
	var v float64
	if err := buf.ReadBigEndian(&v); err != nil {
		return 0, err
	}
	return v, nil
}

func (buf *Buffer) WriteFloat64BE(v float64) error {
	return buf.WriteBigEndian(v)
}
