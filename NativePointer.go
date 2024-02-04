package ffi

/*
#include <string.h>
#include <stdint.h>
*/
import "C"
import (
	"fmt"
	"reflect"
	"sync"
	"unicode/utf16"
	"unsafe"
)

var globalPointer sync.Map

var NULL = Ptr(0)

type NativePointer struct {
	ptr unsafe.Pointer
}
type Numeric interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64 | uintptr | unsafe.Pointer
}

func Ptr[T Numeric](v T) NativePointer {
	vl := reflect.ValueOf(v)
	if vl.Kind() == reflect.Uintptr {
		return NativePointer{
			ptr: unsafe.Pointer(vl.Interface().(uintptr)),
		}
	}
	if vl.Kind() == reflect.UnsafePointer {
		return NativePointer{
			ptr: vl.UnsafePointer(),
		}
	}
	return NativePointer{
		ptr: unsafe.Pointer(uintptr(reflect.ValueOf(v).Int())),
	}
}
func (p NativePointer) Ptr() unsafe.Pointer {
	return p.ptr
}
func (p NativePointer) String() string {
	return p.ToString()
}
func (p NativePointer) ToString() string {
	return fmt.Sprintf("%p", p.ptr)
}
func (p NativePointer) Sub(n int) NativePointer {
	return Ptr(unsafe.Add(p.ptr, -n))
}
func (p NativePointer) Add(n int) NativePointer {
	return Ptr(unsafe.Add(p.ptr, n))
}
func (p NativePointer) ToUinptr() uintptr {
	return uintptr(p.ptr)
}

func (p NativePointer) ReadCString(n int) string {
	str := C.GoString((*C.char)(p.ptr))
	if len(str) == 0 {
		return ""
	}
	if n < 0 {
		return str
	}
	return str[:n]
}
func (p NativePointer) ReadUtf8String(n int) string {
	buf := make([]byte, n)
	copyBytes(p.ptr, n)
	s := unsafe.String(unsafe.SliceData(buf), n)
	if len(s) == 0 {
		return ""
	}
	if n < 0 {
		return s
	}
	return s[:n]
}
func (p NativePointer) ReadInt() int {
	return int(*(*C.int)(p.ptr))
}
func (p NativePointer) ReadFloat() float32 {
	return float32(*(*C.float)(p.ptr))
}
func (p NativePointer) ReadDouble() float64 {
	return float64(*(*C.double)(p.ptr))
}
func (p NativePointer) ReadByteArray(n int) []byte {
	return unsafe.Slice((*byte)(p.ptr), n)
}
func (p NativePointer) ReadLong() int64 {
	return int64(*(*C.long)(p.ptr))
}
func (p NativePointer) ReadPointer() NativePointer {
	return Ptr(unsafe.Pointer(unsafe.Pointer(*(**C.void)(p.ptr))))
}
func (p NativePointer) ReadS8() int8 {
	return int8(*(*C.int8_t)(p.ptr))
}
func (p NativePointer) ReadS16() int16 {
	return int16(*(*C.int16_t)(p.ptr))
}
func (p NativePointer) ReadS32() int32 {
	return int32(*(*C.int32_t)(p.ptr))
}
func (p NativePointer) ReadS64() int64 {
	return int64(*(*C.int64_t)(p.ptr))
}
func (p NativePointer) ReadShort() int16 {
	return int16(*(*C.int16_t)(p.ptr))
}
func (p NativePointer) ReadU8() uint8 {
	return uint8(*(*C.uint8_t)(p.ptr))
}
func (p NativePointer) ReadU16() uint16 {
	return uint16(*(*C.uint16_t)(p.ptr))
}
func (p NativePointer) ReadU32() uint32 {
	return uint32(*(*C.uint32_t)(p.ptr))
}
func (p NativePointer) ReadU64() uint64 {
	return uint64(*(*C.uint64_t)(p.ptr))
}
func (p NativePointer) ReadUShort() uint16 {
	return uint16(*(*C.ushort)(p.ptr))
}
func (p NativePointer) ReadUint() uint {
	return uint(*(*C.uint)(p.ptr))
}
func (p NativePointer) ReadULong() uint64 {
	return uint64(*(*C.ulong)(p.ptr))
}
func (p NativePointer) ReadUTF16String(n int) string {
	bt := unsafe.Slice((*uint16)(p.ptr), n)
	return string(utf16.Decode(bt))
}
func (p NativePointer) IsNull() bool {
	return p.ToUinptr() == 0
}

func (p NativePointer) WritePointer(d NativePointer) NativePointer {
	*(**C.void)(p.ptr) = (*C.void)(d.ptr)
	return p
}
func (p NativePointer) WriteDouble(d float64) NativePointer {
	*(*C.double)(p.ptr) = C.double(d)
	return p
}
func (p NativePointer) WriteFloat(d float32) NativePointer {
	*(*C.float)(p.ptr) = C.float(d)
	return p
}
func (p NativePointer) WriteShort(d int16) NativePointer {
	*(*C.short)(p.ptr) = C.short(d)
	return p
}
func (p NativePointer) WriteInt(d int) NativePointer {
	*(*C.int)(p.ptr) = C.int(d)
	return p
}
func (p NativePointer) WriteLong(d int64) NativePointer {
	*(*C.long)(p.ptr) = C.long(d)
	return p
}
func (p NativePointer) WriteUShort(d uint16) NativePointer {
	*(*C.uint16_t)(p.ptr) = C.uint16_t(d)
	return p
}
func (p NativePointer) WriteUInt(d uint) NativePointer {
	*(*C.int32_t)(p.ptr) = C.int32_t(d)
	return p
}
func (p NativePointer) WriteULong(d uint64) NativePointer {
	*(*C.int64_t)(p.ptr) = C.int64_t(d)
	return p
}
func (p NativePointer) WriteS8(d int8) NativePointer {
	*(*C.int8_t)(p.ptr) = C.int8_t(d)
	return p
}
func (p NativePointer) WriteS16(d int16) NativePointer {
	*(*C.int16_t)(p.ptr) = C.int16_t(d)
	return p
}
func (p NativePointer) WriteS32(d int32) NativePointer {
	*(*C.int32_t)(p.ptr) = C.int32_t(d)
	return p
}
func (p NativePointer) WriteS64(d int64) NativePointer {
	*(*C.int64_t)(p.ptr) = C.int64_t(d)
	return p
}
func (p NativePointer) WriteU8(d uint8) NativePointer {
	*(*C.uint8_t)(p.ptr) = C.uint8_t(d)
	return p
}
func (p NativePointer) WriteU16(d uint16) NativePointer {
	*(*C.uint16_t)(p.ptr) = C.uint16_t(d)
	return p
}
func (p NativePointer) WriteU32(d uint32) NativePointer {
	*(*C.uint32_t)(p.ptr) = C.uint32_t(d)
	return p
}
func (p NativePointer) WriteU64(d uint64) NativePointer {
	*(*C.uint64_t)(p.ptr) = C.uint64_t(d)
	return p
}
func (p NativePointer) WriteByteArray(b []byte) NativePointer {
	for i := 0; i < len(b); i++ {
		*(*byte)(unsafe.Pointer(uintptr(p.ptr) + uintptr(i))) = b[i]
	}
	return p
}
func (p NativePointer) WriteUtf8String(s string) NativePointer {
	for i := 0; i < len(s); i++ {
		*(*C.char)(unsafe.Pointer(uintptr(p.ptr) + uintptr(i))) = C.char(s[i])
	}
	return p
}
func (p NativePointer) WriteUtf16String(s string) NativePointer {
	for i := 0; i < len(s); i++ {
		*(*rune)(unsafe.Pointer(uintptr(p.ptr) + uintptr(i))) = rune(s[i])
	}
	return p
}

//
//func (p NativePointer) Dump(n int) string {
//	return hex.Dump(p.ReadByteArray(n))
//}
//func (p NativePointer) Pointer() unsafe.Pointer {
//	return p.ptr
//}
//func (p NativePointer) ReadCString(n int) string {
//	s := GoString((*C.char)(p.ptr))
//	//if n > 1 {
//	//	s = s[:n]
//	//}
//	return s
//}
//func (p NativePointer) ReadUtf8String(n int) string {
//	s := unsafe.String(unsafe.SliceData(p.ReadByteArray(n)), n)
//	return s
//}
//func (p NativePointer) ReadByteArray(n int) []byte {
//	bt := unsafe.Slice((*byte)(p.ptr), n)
//	return bt
//}
//func (p NativePointer) ReadPointer() NativePointer {
//	p1 := (**C.void)(p.ptr)
//	p2 := unsafe.Pointer(*p1)
//	return Ptr(p2)
//}
//
//func (p NativePointer) ReadUint() uint {
//	return p.ReadPointer().ToUint()
//}
//func (p NativePointer) ReadInt8() int8 {
//	return int8(p.ReadPointer().ToInt8())
//}
//func (p NativePointer) ReadUint8() uint8 {
//	return uint8(p.ReadPointer().ToUint8())
//}
//func (p NativePointer) ReadInt16() int16 {
//	return int16(p.ReadPointer().ToInt16())
//}
//func (p NativePointer) ReadUint16() uint16 {
//	return uint16(p.ReadPointer().ToUint16())
//}
//func (p NativePointer) ReadInt32() int32 {
//	return int32(p.ReadPointer().ToInt32())
//}
//func (p NativePointer) ReadUint32() uint32 {
//	return uint32(p.ReadPointer().ToUint32())
//}
//func (p NativePointer) ReadInt64() int64 {
//	return int64(p.ReadPointer().ToInt64())
//}
//func (p NativePointer) ReadUint64() uint64 {
//	return uint64(p.ReadPointer().ToUint64())
//}
//func (p NativePointer) ReadFloat() float32 {
//	p1 := (*C.float)(p.ptr)
//	return float32(*p1)
//}
//func (p NativePointer) ReadDouble() float64 {
//	p1 := (*C.double)(p.ptr)
//	return float64(*p1)
//}

//func (p NativePointer) WriteInt(n int) {
//	np := unsafe.Pointer(p.ReadPointer().ptr)
//	*(*C.int)(np) = C.int(n)
//}
//func (p NativePointer) WriteUint(n uint) {
//	np := unsafe.Pointer(p.ReadPointer().ptr)
//	*(*C.uint)(np) = C.uint(n)
//}
//func (p NativePointer) WriteInt8(n int8) {
//	np := unsafe.Pointer(p.ReadPointer().ptr)
//	*(*C.int8_t)(np) = C.int8_t(n)
//}
//func (p NativePointer) WriteUint8(n uint8) {
//	np := unsafe.Pointer(p.ReadPointer().ptr)
//	*(*C.uint8_t)(np) = C.uint8_t(n)
//}
//func (p NativePointer) WriteInt16(n int16) {
//	np := unsafe.Pointer(p.ReadPointer().ptr)
//	*(*C.int16_t)(np) = C.int16_t(n)
//}
//func (p NativePointer) WriteUint16(n uint16) {
//	np := unsafe.Pointer(p.ReadPointer().ptr)
//	*(*C.uint16_t)(np) = C.uint16_t(n)
//}
//func (p NativePointer) WriteInt32(n int32) {
//	np := unsafe.Pointer(p.ReadPointer().ptr)
//	*(*C.int32_t)(np) = C.int32_t(n)
//}
//func (p NativePointer) WriteUint32(n uint32) {
//	np := unsafe.Pointer(p.ReadPointer().ptr)
//	*(*C.uint32_t)(np) = C.uint32_t(n)
//}
//func (p NativePointer) WriteInt64(n int64) {
//	np := unsafe.Pointer(p.ReadPointer().ptr)
//	*(*C.int64_t)(np) = C.int64_t(n)
//}
//func (p NativePointer) WriteUint64(n uint64) {
//	np := unsafe.Pointer(p.ReadPointer().ptr)
//	*(*C.uint64_t)(np) = C.uint64_t(n)
//}
//func (p NativePointer) WritePointer(n NativePointer) {
//	np := unsafe.Pointer(p.ReadPointer().ptr)
//	*((**C.void)(np)) = (*C.void)(unsafe.Pointer(n.ptr))
//}
//func (p NativePointer) WriteUtf8String(s string) {
//	Memcpy(unsafe.Pointer(p.ptr), unsafe.Pointer(unsafe.StringData(s)), int64(len(s)))
//}
//func (p NativePointer) WriteByteArray(s []byte) {
//	Memcpy(unsafe.Pointer(p.ptr), unsafe.Pointer(&s[0]), int64(len(s)))
//}

func allocGlobalUtf8String() NativePointer {
	p := allocGlobalUtf8String()
	globalPointer.Store(&p, p)
	return p
}
func allocUtf8String(s string) NativePointer {
	b := make([]byte, len(s)+1)
	copy(b, s)
	b[len(b)-1] = 0x00
	return Ptr(unsafe.Pointer(&b[0]))
}
