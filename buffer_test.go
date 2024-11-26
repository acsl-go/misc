package misc

import "testing"

func TestBuffer(t *testing.T) {
	buf := NewBufferPool(1024).Get()
	buf.AddRef()
	buf.Release()
	buf.WriteLittleEndian(int64(-1))
	buf.Seek(0, 0)
	var v int64
	buf.ReadLittleEndian(&v)
	if v != -1 {
		t.Fatal(v)
	}
	print(v)
	buf.Release()
}
