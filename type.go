package ffi

/*
#include <ffi.h>
#include <stdint.h>
*/
import "C"
import (
	"errors"
	"github.com/samber/lo"
	"reflect"
	"unsafe"
)

type FFIType = C.ffi_type
type RetTypeName = string
type ArgTypeName = string

type NativeABI = C.ffi_abi

const (
	FirstAbi   NativeABI = C.FFI_FIRST_ABI
	DefaultAbi NativeABI = C.FFI_DEFAULT_ABI
	LastAbi    NativeABI = C.FFI_LAST_ABI
)

var (
	Tint    = "int"
	TUint   = "uint"
	TLong   = "long"
	TULong  = "ulong"
	TChar   = "char"
	TUChar  = "uchar"
	TFloat  = "float"
	TDouble = "double"
	Tint8   = "int8"
	TUint8  = "uint8"
	Tint16  = "int16"
	TUint16 = "uint16"
	Tint32  = "int32"
	TUint32 = "uint32"
	TBool   = "bool"
)
var (
	TVoid    = "void"
	TPointer = "pointer"
	TSizeT   = "uint64"
	TSSizeT  = "int64"
	TInt64   = "int64"
	TUint64  = "uint64"
)

func ConvertStringTypeToFFIType(tpName string) *FFIType {
	switch tpName {
	case Tint:
		return &C.ffi_type_sint32
	case TUint:
		return &C.ffi_type_uint32
	case TLong:
		return &C.ffi_type_sint64
	case TULong:
		return &C.ffi_type_uint64
	case TChar:
		return &C.ffi_type_sint8
	case TUChar:
		return &C.ffi_type_uint8
	case TFloat:
		return &C.ffi_type_float
	case TDouble:
		return &C.ffi_type_double
	case Tint8:
		return &C.ffi_type_sint8
	case TUint8:
		return &C.ffi_type_uint8
	case Tint16:
		return &C.ffi_type_sint16
	case TUint32:
		return &C.ffi_type_uint32
	case TBool:
		return &C.ffi_type_sint8
	case TVoid:
		return &C.ffi_type_void
	case TPointer:
		return &C.ffi_type_pointer
	case TSizeT:
		return &C.ffi_type_uint64
	case TSSizeT:
		return &C.ffi_type_sint64
	case TInt64:
		return &C.ffi_type_sint64
	case TUint64:
		return &C.ffi_type_uint64
	}
	panic(errors.New("convert error"))
}
func ConvertAnyToFFIValue(tpName string, v any) unsafe.Pointer {
	vl := reflect.ValueOf(v)
	switch tpName {
	case Tint:
		return unsafe.Pointer(lo.ToPtr(C.int(vl.Int())))
	case TUint:
		return unsafe.Pointer(lo.ToPtr(C.uint(vl.Int())))
	case TLong:
		return unsafe.Pointer(lo.ToPtr(C.long(vl.Int())))
	case TULong:
		return unsafe.Pointer(lo.ToPtr(C.ulong(vl.Int())))
	case TChar:
		return unsafe.Pointer(lo.ToPtr(C.char(vl.Int())))
	case TUChar:
		return unsafe.Pointer(lo.ToPtr(C.uchar(vl.Int())))
	case TFloat:
		if vl.CanFloat() {
			return unsafe.Pointer(lo.ToPtr(C.float(vl.Float())))
		} else {
			return unsafe.Pointer(lo.ToPtr(C.float(vl.Int())))
		}
	case TDouble:
		if vl.CanFloat() {
			return unsafe.Pointer(lo.ToPtr(C.double(vl.Float())))
		} else {
			return unsafe.Pointer(lo.ToPtr(C.float(vl.Int())))
		}
	case Tint8:
		return unsafe.Pointer(lo.ToPtr(C.int8_t(vl.Int())))
	case TUint8:
		return unsafe.Pointer(lo.ToPtr(C.uint8_t(vl.Uint())))
	case Tint16:
		return unsafe.Pointer(lo.ToPtr(C.int16_t(vl.Int())))
	case TUint32:
		return unsafe.Pointer(lo.ToPtr(C.uint32_t(vl.Uint())))
	case TBool:
		return unsafe.Pointer(lo.ToPtr(C.int(vl.Int())))
	case TPointer:
		switch vl.Interface().(type) {
		case NativePointer:
			return unsafe.Pointer(lo.ToPtr(unsafe.Pointer(vl.Interface().(NativePointer).Ptr())))
		case unsafe.Pointer:
			return unsafe.Pointer(lo.ToPtr(unsafe.Pointer(vl.Interface().(unsafe.Pointer))))
		case uintptr:
			return unsafe.Pointer(lo.ToPtr(unsafe.Pointer(vl.Interface().(uintptr))))
		default:
			return unsafe.Pointer(lo.ToPtr(unsafe.Pointer(vl.Pointer())))
		}
	case TSizeT:
		return unsafe.Pointer(lo.ToPtr(C.uint64_t(vl.Uint())))
	case TSSizeT:
		return unsafe.Pointer(lo.ToPtr(C.int64_t(vl.Int())))
	case TInt64:
		return unsafe.Pointer(lo.ToPtr(C.int64_t(vl.Int())))
	case TUint64:
		return unsafe.Pointer(lo.ToPtr(C.uint64_t(vl.Uint())))
	}
	panic(errors.New("convert error"))
}
func ConvertFFIValueToAny(tpName string, pointer unsafe.Pointer) any {
	valPtr := Ptr(pointer)
	switch tpName {
	case Tint:
		return valPtr.ReadInt()
	case TUint:
		return valPtr.ReadUint()
	case TLong:
		return valPtr.ReadLong()
	case TULong:
		return valPtr.ReadULong()
	case TChar:
		return valPtr.ReadS8()
	case TUChar:
		return valPtr.ReadU8()
	case TFloat:
		return valPtr.ReadFloat()
	case TDouble:
		return valPtr.ReadFloat()
	case Tint8:
		return valPtr.ReadS8()
	case TUint8:
		return valPtr.ReadU8()
	case Tint16:
		return valPtr.ReadS16()
	case TUint16:
		return valPtr.ReadU16()
	case Tint32:
		return valPtr.ReadS32()
	case TUint32:
		return valPtr.ReadU32()
	case TBool:
		return valPtr.ReadS8()
	case TPointer:
		return valPtr.ReadPointer()
	case TSizeT:
		return valPtr.ReadU64()
	case TSSizeT:
		return valPtr.ReadS64()
	case TInt64:
		return valPtr.ReadS64()
	case TUint64:
		return valPtr.ReadU64()
	case TVoid:
		return Ptr(0)
	}
	panic(errors.New("convert error"))
}
func ConvertFFIRetValue(tpName string, pointer unsafe.Pointer) NativePointer {
	valPtr := Ptr(pointer)
	switch tpName {
	case Tint:
		return Ptr(valPtr.ReadInt())
	case TUint:
		return Ptr(valPtr.ReadUint())
	case TLong:
		return Ptr(valPtr.ReadLong())
	case TULong:
		return Ptr(valPtr.ReadULong())
	case TChar:
		return Ptr(valPtr.ReadS8())
	case TUChar:
		return Ptr(valPtr.ReadU8())
	case TFloat:
		return valPtr
	case TDouble:
		return valPtr
	case Tint8:
		return Ptr(valPtr.ReadS8())
	case TUint8:
		return Ptr(valPtr.ReadU8())
	case Tint16:
		return Ptr(valPtr.ReadS16())
	case TUint16:
		return Ptr(valPtr.ReadU16())
	case Tint32:
		return Ptr(valPtr.ReadS32())
	case TUint32:
		return Ptr(valPtr.ReadU32())
	case TBool:
		return Ptr(valPtr.ReadS8())
	case TPointer:
		return valPtr.ReadPointer()
	case TSizeT:
		return Ptr(valPtr.ReadU64())
	case TSSizeT:
		return Ptr(valPtr.ReadS64())
	case TInt64:
		return Ptr(valPtr.ReadS64())
	case TUint64:
		return Ptr(valPtr.ReadU64())
	case TVoid:
		return Ptr(0)
	}
	panic(errors.New("convert error"))
}
func WriteRetValue(retPtr NativePointer, tpName string, v any) {
	if tpName == TVoid {
		return
	}
	vl := reflect.ValueOf(v)
	switch tpName {
	case Tint:
		retPtr.WriteInt(int(vl.Int()))
	case TUint:
		retPtr.WriteUInt(uint(vl.Uint()))
	case TLong:
		retPtr.WriteLong(vl.Int())
	case TULong:
		retPtr.WriteULong(vl.Uint())
	case TChar:
		retPtr.WriteS8(int8(vl.Int()))
	case TUChar:
		retPtr.WriteU8(uint8(vl.Uint()))
	case TFloat:
		retPtr.WriteFloat(float32(vl.Float()))
	case TDouble:
		retPtr.WriteDouble(float64(vl.Float()))
	case Tint8:
		retPtr.WriteS8(int8(vl.Int()))
	case TUint8:
		retPtr.WriteU8(uint8(vl.Uint()))
	case Tint16:
		retPtr.WriteS16(int16(vl.Int()))
	case TUint16:
		retPtr.WriteU16(uint16(vl.Int()))
	case Tint32:
		retPtr.WriteS32(int32(vl.Int()))
	case TUint32:
		retPtr.WriteU32(uint32(vl.Int()))
	case TBool:
		retPtr.WriteS8(int8(vl.Int()))
	case TPointer:
		retPtr.WritePointer(vl.Interface().(NativePointer))
	case TSizeT:
		retPtr.WriteU64(uint64(vl.Int()))
	case TSSizeT:
		retPtr.WriteS64(int64(vl.Int()))
	case TInt64:
		retPtr.WriteS64(int64(vl.Int()))
	case TUint64:
		retPtr.WriteU64(uint64(vl.Int()))
	default:
		panic(errors.New("convert error"))
	}

}
func MakeFFIRetValue(tpName string) unsafe.Pointer {
	switch tpName {
	case Tint:
		return unsafe.Pointer(lo.ToPtr(C.int(0)))
	case TUint:
		return unsafe.Pointer(lo.ToPtr(C.uint(0)))
	case TLong:
		return unsafe.Pointer(lo.ToPtr(C.long(0)))
	case TULong:
		return unsafe.Pointer(lo.ToPtr(C.ulong(0)))
	case TChar:
		return unsafe.Pointer(lo.ToPtr(C.char(0)))
	case TUChar:
		return unsafe.Pointer(lo.ToPtr(C.uchar(0)))
	case TFloat:
		return unsafe.Pointer(lo.ToPtr(C.float(0)))
	case TDouble:
		return unsafe.Pointer(lo.ToPtr(C.double(0)))
	case Tint8:
		return unsafe.Pointer(lo.ToPtr(C.int8_t(0)))
	case TUint8:
		return unsafe.Pointer(lo.ToPtr(C.uint8_t(0)))
	case Tint16:
		return unsafe.Pointer(lo.ToPtr(C.int16_t(0)))
	case TUint32:
		return unsafe.Pointer(lo.ToPtr(C.uint32_t(0)))
	case TBool:
		return unsafe.Pointer(lo.ToPtr(C.int(0)))
	case TPointer:
		return unsafe.Pointer(lo.ToPtr((*C.void)(unsafe.Pointer(uintptr(0)))))
	case TSizeT:
		return unsafe.Pointer(lo.ToPtr(C.uint64_t(0)))
	case TSSizeT:
		return unsafe.Pointer(lo.ToPtr(C.int64_t(0)))
	case TInt64:
		return unsafe.Pointer(lo.ToPtr(C.int64_t(0)))
	case TUint64:
		return unsafe.Pointer(lo.ToPtr(C.uint64_t(0)))
	}
	panic(errors.New("convert error"))
}
