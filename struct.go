package unsafe2

import (
	"reflect"
	"unsafe"
)

type Struct struct {
	// Fields is a slice containing the sizes of all the fields
	Fields []uintptr
	// Data is the raw data of the struct
	Data []byte
}

// AllocStruct allocates a struct object with the appropriate size
func AllocStruct(fields []uintptr) *Struct {
	var size uintptr
	for _, f := range fields {
		size += f
	}

	return &Struct{
		Fields: fields, Data: make([]byte, size),
	}
}

// NewStructFromPointer makes a struct with the bytes of the pointer at ptr with the specified fields
func NewStructFromPointer(ptr unsafe.Pointer, size uintptr, fields []uintptr) *Struct {
	return &Struct{
		Fields: fields,
		Data:   PtrBytes(ptr, size),
	}
}

// NewStructFromStruct makes a struct with the bytes of the value of data with the specified fields
func NewStructFromStruct[T any](data *T, fields []uintptr) *Struct {
	ptr := unsafe.Pointer(data)
	size := unsafe.Sizeof(*data)

	return NewStructFromPointer(ptr, size, fields)
}

// NewStructFromStructCopy makes a struct with data with the specified fields
func NewStructFromStructCopy(data any, fields []uintptr) *Struct {
	a := InterfaceData(data)
	typ := (*Type)(unsafe.Pointer(a.Type))
	size := typ.Size

	ptr := unsafe.Pointer(a.Value)

	return NewStructFromPointer(ptr, size, fields)
}

// NewStructReflect makes a struct with the bytes of the value of data and finds its fields by reflecting
func NewStructReflect[T any](data *T) *Struct {
	ptr := unsafe.Pointer(data)

	t := reflect.ValueOf(*data).Type()
	fields := make([]uintptr, t.NumField())
	for i := range fields {
		fields[i] = t.Field(i).Type.Size()
	}

	return NewStructFromPointer(ptr, t.Size(), fields)
}

// NewStructReflectCopy makes a struct with the value of data and finds its fields by reflecting
func NewStructReflectCopy(data any) *Struct {
	ptr := unsafe.Pointer(InterfaceData(data).Value)

	t := reflect.ValueOf(data).Type()
	fields := make([]uintptr, t.NumField())
	for i := range fields {
		fields[i] = t.Field(i).Type.Size()
	}

	return NewStructFromPointer(ptr, t.Size(), fields)
}

func (s *Struct) off(i int) uintptr {
	if i == 0 {
		return 0
	}
	var off uintptr
	for _, f := range s.Fields[:i] {
		off += f
	}

	return off
}

func (s *Struct) SetField(i int, v []byte) {
	size := s.Fields[i]
	off := s.off(i)

	copy(s.Data[off:off+size], v)
}

func (s *Struct) SetFieldPtr(i int, ptr unsafe.Pointer) {
	size := s.Fields[i]
	off := s.off(i)

	v := PtrBytes(ptr, size)

	copy(s.Data[off:off+size], v)
}

func (s *Struct) SetFieldData(i int, data any) {
	size := s.Fields[i]
	off := s.off(i)

	ptr := unsafe.Pointer(InterfaceData(data).Value)

	v := PtrBytes(ptr, size)
	copy(s.Data[off:off+size], v)
}

func (s *Struct) Field(i int) []byte {
	size := s.Fields[i]
	off := s.off(i)

	return s.Data[off : off+size]
}

func (s *Struct) FieldPtr(i int) unsafe.Pointer {
	size := s.Fields[i]
	off := s.off(i)

	return unsafe.Pointer(unsafe.SliceData(s.Data[off : off+size]))
}

func StructFieldPtrCast[T any](s *Struct, i int) *T {
	return (*T)(s.FieldPtr(i))
}

func (s *Struct) CopyField(i int, dst unsafe.Pointer) {
	size := s.Fields[i]
	off := s.off(i)

	dstslice := unsafe.Slice((*byte)(dst), size)

	copy(dstslice, s.Data[off:off+size])
}
