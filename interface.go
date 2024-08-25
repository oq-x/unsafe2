package github.com/oq-x/unsafe2

import (
	"unsafe"
)

var (
	IntType   = TypeOf(0)
	Int8Type  = TypeOf(int8(0))
	Int16Type = TypeOf(int16(0))
	Int32Type = TypeOf(int32(0))
	Int64Type = TypeOf(int64(0))

	UintType    = TypeOf(uint(0))
	Uint8Type   = TypeOf(uint8(0))
	Uint16Type  = TypeOf(uint16(0))
	Uint32Type  = TypeOf(uint32(0))
	Uint64Type  = TypeOf(uint64(0))
	UintptrType = TypeOf(uintptr(0))

	Float32Type = TypeOf(float32(0))
	Float64Type = TypeOf(0.0)

	Complex64Type  = TypeOf(complex64(0))
	Complex128Type = TypeOf(complex128(0))

	StringType = TypeOf("")
)

type Type struct {
	Size       uintptr
	PtrBytes   uintptr
	Hash       uint32
	TFlag      uint8
	Align      uint8
	FieldAlign uint8

	Kind uint8

	Equal     func(unsafe.Pointer, unsafe.Pointer) bool
	GCData    *byte
	Name      int32
	PtrToThis int32
}

func NewAny(t *Type, v uintptr) *any {
	return (*any)(unsafe.Pointer(&Interface{Type: uintptr(unsafe.Pointer(t)), Value: v}))
}

type Interface struct {
	Type  uintptr
	Value uintptr
}

func TypeOf(a any) *Type {
	return (*Type)(unsafe.Pointer(InterfaceData(a).Type))
}

func ValueOf(a any) unsafe.Pointer {
	return unsafe.Pointer(InterfaceData(a).Value)
}

func InterfaceData(a any) *Interface {
	return (*Interface)(unsafe.Pointer(&a))
}

// InterfaceCastCopy copies the concrete value of a into v. size should match the actual concrete's type value
func InterfaceCastCopy(a any, v unsafe.Pointer, size uintptr) {
	data := InterfaceData(a)

	src := PtrBytes(unsafe.Pointer(data.Value), size)
	dst := PtrBytes(v, size)
	copy(dst, src)
}

// InterfaceCast unsafely casts the concrete value of a to the type specified.
func InterfaceCast[T any](a any) *T {
	data := InterfaceData(a)
	return (*T)(unsafe.Pointer(data.Value))
}

// InterfaceSetPtr sets the concrete value of a to the value behind v. The sizes must match
func InterfaceSetPtr(a any, v unsafe.Pointer, size uintptr) {
	dst := PtrBytes(unsafe.Pointer(InterfaceData(a).Value), size)
	src := PtrBytes(v, size)
	copy(dst, src)
}

// InterfaceSetValue sets the concrete value of a to v. The sizes must match
func InterfaceSetValue[T any](a any, v T) {
	size := unsafe.Sizeof(v)
	dst := PtrBytes(unsafe.Pointer(InterfaceData(a).Value), size)
	src := PtrBytes(unsafe.Pointer(&v), size)
	copy(dst, src)
}

// Comparable checks if the type of a is comparable
func Comparable(a any) bool {
	return TypeOf(a).Equal != nil
}

// Equal uses the equal function of a's type to check if a and b are equal
func Equal(a, b any) bool {
	typ := TypeOf(a)
	ai := unsafe.Pointer(InterfaceData(a).Value)
	bi := unsafe.Pointer(InterfaceData(b).Value)

	return typ.Equal(ai, bi)
}
