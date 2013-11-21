// Package tiger implements the Tiger hash algorithm
// https://github.com/cxmcc/tiger
package tiger

import "hash"

// The size of a Tiger hash value in bytes
const Size = 24

// The blocksize of Tiger hash function in bytes
const BlockSize = 64

const (
	chunk = 64
	initA = 0x0123456789abcdef
	initB = 0xfedcba9876543210
	initC = 0xf096a5b4c3b2e187
)

type digest struct {
	a      uint64
	b      uint64
	c      uint64
	x      [chunk]byte
	nx     int
	length uint64
	ver    int
}

func (d *digest) Reset() {
	d.a = initA
	d.b = initB
	d.c = initC
	d.nx = 0
	d.length = 0
}

// New returns a new hash.Hash computing the Tiger hash value
func New() hash.Hash {
	d := new(digest)
	d.Reset()
	d.ver = 1
	return d
}

// New returns a new hash.Hash computing the Tiger2 hash value
func New2() hash.Hash {
	d := new(digest)
	d.Reset()
	d.ver = 2
	return d
}

func (d *digest) BlockSize() int {
	return BlockSize
}

func (d *digest) Size() int {
	return Size
}

func (d *digest) Write(p []byte) (nn int, err error) {
	nn = len(p)
	d.length += uint64(nn)
	if d.nx > 0 {
		n := len(p)
		if n > chunk-d.nx {
			n = chunk - d.nx
		}
		for i := 0; i < n; i++ {
			d.x[d.nx+i] = p[i]
		}
		d.nx += n
		if d.nx == chunk {
			compress(d, d.x[0:chunk])
			d.nx = 0
		}
		p = p[n:]
	}
	if len(p) >= chunk {
		n := len(p) &^ (chunk - 1)
		compress(d, p[:n])
		p = p[n:]
	}
	if len(p) > 0 {
		d.nx = copy(d.x[:], p)
	}
	return
}

func (d0 *digest) Sum(in []byte) []byte {
	d := *d0

	length := d.length
	var tmp [64]byte
	if d.ver == 1 {
		tmp[0] = 0x01
	} else {
		tmp[0] = 0x80

	}

	if length&0x3f < 56 {
		d.Write(tmp[0 : 56-length&0x3f])
	} else {
		d.Write(tmp[0 : 64+56-length&0x3f])
	}

	length <<= 3
	for i := uint(0); i < 8; i++ {
		tmp[i] = byte(length >> (8 * i))
	}
	d.Write(tmp[0:8])

	if d.nx != 0 {
		panic("d.nx != 0")
	}

	for i := uint(0); i < 8; i++ {
		tmp[i] = byte(d.a >> (8 * i))
		tmp[i+8] = byte(d.b >> (8 * i))
		tmp[i+16] = byte(d.c >> (8 * i))
	}

	return append(in, tmp[:24]...)
}
