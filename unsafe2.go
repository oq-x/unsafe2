// Very unsafe operations
package github.com/oq-x/unsafe2

import (
	"unsafe"
)

// PtrBytes returns the byte representation of the value behind ptr. Any value changed in the slice, will also change in the original value
func PtrBytes(ptr unsafe.Pointer, size uintptr) []byte {
	return unsafe.Slice((*byte)(ptr), size)
}

// DataBytes returns the byte representation of the data. Any value changed in the slice, will also change in the original value
func DataBytes[T any](data *T) []byte {
	return unsafe.Slice((*byte)(unsafe.Pointer(data)), unsafe.Sizeof(*data))
}

// BytesCopy sets the value of dst to data
func BytesCopy(data []byte, dst unsafe.Pointer) {
	copy(PtrBytes(dst, uintptr(len(data))), data)
}
