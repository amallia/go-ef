package ef

import (
	"errors"
	"github.com/willf/bitset"
	"log"
	"math"
)

const (
	EFInfo = `Universe: %d
Elements: %d
Lower_bits: %d
Higher_bits_length: %d
Mask: 0b%b
Lower_bits offset: %d
Bitvector length: %d
`
)

type EliasFano struct {
	universe           uint64
	n                  uint64
	lower_bits         uint64
	higher_bits_length uint64
	mask               uint64
	lower_bits_offset  uint64
	bv_len             uint64
	b                  *bitset.BitSet
	cur_value          uint64
	position           uint64
	high_bits_pos      uint64
}

func New(universe uint64, n uint64) *EliasFano {
	var lower_bits uint64
	if lower_bits = 0; universe > n {
		lower_bits = msb(universe / n)
	}
	higher_bits_length := n + (universe >> lower_bits) + 2
	mask := (uint64(1) << lower_bits) - 1
	lower_bits_offset := higher_bits_length
	bv_len := lower_bits_offset + n*uint64(lower_bits)
	b := bitset.New(uint(bv_len))
	return &EliasFano{universe, n, lower_bits, higher_bits_length, mask, lower_bits_offset, bv_len, b, 0, 0, 0}
}

func (ef *EliasFano) Compress(elems []uint64) {
	last := uint64(0)

	for i, elem := range elems {
		if i > 0 && elem < last {
			log.Fatal("Sequence is not sorted")
		}
		if elem > ef.universe {
			log.Fatalf("Element %d is greater than universe", elem)
		}
		high := (elem >> ef.lower_bits) + uint64(i) + 1
		low := elem & ef.mask
		ef.b.Set(uint(high))
		offset := ef.lower_bits_offset + uint64(i)*ef.lower_bits
		SetBits(ef.b, offset, low, ef.lower_bits)
		last = elem
		if i == 0 {
			ef.cur_value = elem
			ef.high_bits_pos = high
		}
	}
}

func SetBits(b *bitset.BitSet, offset uint64, bits uint64, length uint64) {
	// TODO: Store reversed
	for i := uint64(0); i < length; i++ {
		val := bits & (1 << (length - i - 1))
		b.SetTo(uint(offset+i+1), val > 0)
	}

}

func (ef *EliasFano) Next() (uint64, error) {
	ef.position++
	if ef.position >= ef.Size() {
		return 0, errors.New("End reached")
	}
	pos := uint(ef.high_bits_pos)
	if pos > 0 {
		pos++
	}
	pos, _ = ef.b.NextSet(pos)
	ef.high_bits_pos = uint64(pos)
	low := uint64(0)
	offset := ef.lower_bits_offset + ef.position*ef.lower_bits
	for i := uint64(0); i < ef.lower_bits; i++ {
		if ef.b.Test(uint(offset + i + 1)) {
			low++
		}
		low = low << 1
	}
	low = low >> 1
	ef.cur_value = uint64(((ef.high_bits_pos - ef.position - 1) << ef.lower_bits) | low)
	return ef.Value(), nil

}

func (ef *EliasFano) Position() uint64 {
	return ef.position
}

func (ef *EliasFano) Reset() {
	ef.position = 0
}

func (ef *EliasFano) Info() {
	log.Printf(EFInfo, ef.universe, ef.n, ef.lower_bits, ef.higher_bits_length, ef.mask, ef.lower_bits_offset, ef.bv_len)
}

func round(a float64) int64 {
	if a < 0 {
		return int64(a - 0.5)
	}
	return int64(a + 0.5)
}

func (ef *EliasFano) Value() uint64 {
	return ef.cur_value
}

func (ef *EliasFano) Size() uint64 {
	return ef.n
}

func (ef *EliasFano) Bitsize() uint64 {
	return uint64(ef.b.BinaryStorageSize())
}

func msb(x uint64) uint64 {
	return uint64(round(math.Log2(float64(x))))
}
