package ef

import (
	"github.com/amallia/go-ef"
	"testing"
)

func TestMembership(t *testing.T) {
	num := uint64(1000)
	obj := ef.New(num, num)
	array := make([]uint64, num)
	for i := range array {
		array[i] = uint64(i)
	}
	obj.Compress(array)
	for _, v := range array {
		obj.Next()
		if obj.Value() != v {
			t.Errorf("%d is not %d. Missing value", obj.Value(), v)
		}
	}
}

// 1001  1010
