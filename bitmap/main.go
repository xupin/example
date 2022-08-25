package main

import (
	"fmt"
)

type BitMap struct {
	bits []byte
	max  uint64
}

func main() {
	bm := NewBitMap(8)
	// fmt.Printf("%08b\n", bm.bits)
	bm.Add(8)
	// fmt.Printf("%08b\n", bm.bits)
	fmt.Println(bm.Lookup(8))
	bm.Del(8)
	// fmt.Printf("%08b\n", bm.bits)
}

func NewBitMap(max uint64) *BitMap {
	bm := &BitMap{}
	bm.max = max
	bm.bits = make([]byte, (max>>3)+1)
	return bm
}

// 将1左移{num % 8}位，按位或运算把对应bit设为1
func (r *BitMap) Add(num uint64) {
	r.bits[num>>3] |= 1 << (num % 8)
}

// 将1左移{num % 8}位后，按位与、异或运算把相同bit设为0
func (r *BitMap) Del(num uint64) {
	r.bits[num>>3] &^= 1 << (num % 8)
}

// 将1左移{num % 8}位后，按位与运算判断是否bit都为1
func (r *BitMap) Lookup(num uint64) bool {
	return r.bits[num>>3]&1<<(num%8) != 0
}
