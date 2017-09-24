package ef

import (
	"testing"
)

func TestMembership(t *testing.T) {
	num := uint64(1000)
	obj := New(num, num)
	array := make([]uint64, num)
	for i := range array {
		array[i] = uint64(i)
	}
	obj.Compress(array)
	for i, v := range array {
		if obj.Value() != v {
			t.Errorf("%d is not %d. Missing value", obj.Value(), v)
		}
		_, err := obj.Next()
		if err != nil {
			if i != len(array)-1 {
				t.Error(err)
			}
		}
	}
}

func TestPosition(t *testing.T) {
	num := uint64(1000)
	obj := New(num, num)
	array := make([]uint64, num)
	for i := range array {
		array[i] = uint64(i)
	}
	obj.Compress(array)
	for i := range array {
		if obj.Position() != uint64(i) {
			t.Errorf("%d is not %d. Wrong position", obj.Position(), i)
		}
		obj.Next()
	}
}

func TestReset(t *testing.T) {
	num := uint64(1000)
	obj := New(num, num)
	array := make([]uint64, num)
	for i := range array {
		array[i] = uint64(i)
	}
	obj.Compress(array)
	if obj.Position() != 0 {
		t.Errorf("Initial position is not 0.")
	}
	obj.Next()
	obj.Reset()
	if obj.Position() != 0 {
		t.Errorf("Position not correctly reset.")
	}
	if obj.Value() != 0 {
		t.Errorf("%d is not %d. Missing value", obj.Value(), 0)
	}
}

func TestMove(t *testing.T) {
	num := uint64(1000)
	obj := New(num, num)
	array := make([]uint64, num)
	for i := range array {
		array[i] = uint64(i)
	}
	obj.Compress(array)
	if obj.Position() != 0 {
		t.Errorf("Initial position is not 0.")
	}

	for i, v := range array {
		obj.Move(uint64(i))
		if obj.Value() != v {
			t.Errorf("%d is not %d. Missing value", obj.Value(), v)
		}
	}
	for i := range array {
		obj.Move(uint64(len(array) - i - 1))
		if obj.Value() != array[len(array)-i-1] {
			t.Errorf("%d is not %d. Missing value", obj.Value(), array[len(array)-i-1])
		}
	}
}
func TestGeneric(t *testing.T) {
	obj := New(1000, 5)
	obj.Compress([]uint64{0, 5, 9, 800, 1000})
	if obj.Value() != 0 {
		t.Errorf("%d is not %d. Missing value", obj.Value(), 0)
	}
	obj.Move(0)
	if obj.Value() != 0 {
		t.Errorf("%d is not %d. Missing value", obj.Value(), 0)
	}
	obj.Move(4)
	if obj.Value() != 1000 {
		t.Errorf("%d is not %d. Missing value", obj.Value(), 1000)
	}
	obj.Reset()
	if obj.Value() != 0 {
		t.Errorf("%d is not %d. Missing value", obj.Value(), 0)
	}
	obj.Next()
	if obj.Value() != 5 {
		t.Errorf("%d is not %d. Missing value", obj.Value(), 5)
	}
	obj.Next()
	if obj.Value() != 9 {
		t.Errorf("%d is not %d. Missing value", obj.Value(), 9)
	}
	obj.Move(1)
	if obj.Value() != 5 {
		t.Errorf("%d is not %d. Missing value", obj.Value(), 5)
	}

}
