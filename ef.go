package ef

import (
	"errors"
	"github.com/willf/bitset"
	"log"
	"math"
)

const (
	efInfo = `Universe: %d
Elements: %d
Lower_bits: %d
Higher_bits_length: %d
Mask: 0b%b
Lower_bits offset: %d
Bitvector length: %d
`
)

// EliasFano codec structure
type EliasFano struct {
	universe         uint64
	n                uint64
	lowerBits        uint64
	higherBitsLength uint64
	mask             uint64
	lowerBitsOffset  uint64
	bvLen            uint64
	b                *bitset.BitSet
	curValue         uint64
	position         uint64
	highBitsPos      uint64
}

// New creates a new empty EliasFano object
func New(universe uint64, n uint64) *EliasFano {
	var lowerBits uint64
	if lowerBits = 0; universe > n {
		lowerBits = msb(universe / n)
	}
	higherBitsLength := n + (universe >> lowerBits) + 2
	mask := (uint64(1) << lowerBits) - 1
	lowerBitsOffset := higherBitsLength
	bvLen := lowerBitsOffset + n*uint64(lowerBits)
	b := bitset.New(uint(bvLen))
	return &EliasFano{universe, n, lowerBits, higherBitsLength, mask, lowerBitsOffset, bvLen, b, 0, 0, 0}
}

// Compress a monotone increasing array of positive integers. It sets the position at the beginning.
func (ef *EliasFano) Compress(elems []uint64) {
	last := uint64(0)

	for i, elem := range elems {
		if i > 0 && elem < last {
			log.Fatal("Sequence is not sorted")
		}
		if elem > ef.universe {
			log.Fatalf("Element %d is greater than universe", elem)
		}
		high := (elem >> ef.lowerBits) + uint64(i) + 1
		low := elem & ef.mask
		ef.b.Set(uint(high))
		offset := ef.lowerBitsOffset + uint64(i)*ef.lowerBits
		setBits(ef.b, offset, low, ef.lowerBits)
		last = elem
		if i == 0 {
			ef.curValue = elem
			ef.highBitsPos = high
		}
	}
}

// Move the internal iterator to the given position and returns a value or an error.
func (ef *EliasFano) Move(position uint64) (uint64, error) {
	if position >= ef.Size() {
		return 0, errors.New("Out of bound")
	}
	if ef.position == position {
		return ef.Value(), nil
	}
	if position < ef.position {
		ef.Reset()
	}
	skip := position - ef.position
	pos := uint(ef.highBitsPos)
	for i := uint64(0); i < skip; i++ {
		pos, _ = ef.b.NextSet(pos + 1)
	}
	ef.highBitsPos = uint64(pos - 1)
	ef.position = position
	ef.readCurrentValue()
	return ef.Value(), nil
}

// Next moves the internal iterator to the next position and returns a value or an error.
func (ef *EliasFano) Next() (uint64, error) {
	ef.position++
	if ef.position >= ef.Size() {
		return 0, errors.New("End reached")
	}
	ef.readCurrentValue()
	return ef.Value(), nil

}

// Position return the current position of the internal iterator.
func (ef *EliasFano) Position() uint64 {
	return ef.position
}

// Reset moves the internal iterator to the beginning.
func (ef *EliasFano) Reset() {
	ef.highBitsPos = 0
	ef.position = 0
	ef.readCurrentValue()
}

// Info prints info regarding the EliasFano codec.
func (ef *EliasFano) Info() {
	log.Printf(efInfo, ef.universe, ef.n, ef.lowerBits, ef.higherBitsLength, ef.mask, ef.lowerBitsOffset, ef.bvLen)
}

// Value returns the value of the current element.
func (ef *EliasFano) Value() uint64 {
	return ef.curValue
}

// Size returns the number of elements encoded.
func (ef *EliasFano) Size() uint64 {
	return ef.n
}

// Bitsize returns the size of the internal bitvector.
func (ef *EliasFano) Bitsize() uint64 {
	return uint64(ef.b.BinaryStorageSize())
}

func setBits(b *bitset.BitSet, offset uint64, bits uint64, length uint64) {
	for i := uint64(0); i < length; i++ {
		val := bits & (1 << (length - i - 1))
		b.SetTo(uint(offset+i+1), val > 0)
	}

}

func (ef *EliasFano) readCurrentValue() {
	pos := uint(ef.highBitsPos)
	if pos > 0 {
		pos++
	}
	pos, _ = ef.b.NextSet(pos)
	ef.highBitsPos = uint64(pos)
	low := uint64(0)
	offset := ef.lowerBitsOffset + ef.position*ef.lowerBits
	for i := uint64(0); i < ef.lowerBits; i++ {
		if ef.b.Test(uint(offset + i + 1)) {
			low++
		}
		low = low << 1
	}
	low = low >> 1
	ef.curValue = uint64(((ef.highBitsPos - ef.position - 1) << ef.lowerBits) | low)
}

func round(a float64) int64 {
	if a < 0 {
		return int64(a - 0.5)
	}
	return int64(a + 0.5)
}

func msb(x uint64) uint64 {
	return uint64(round(math.Log2(float64(x))))
}
