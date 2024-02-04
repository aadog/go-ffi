package ffi

/*
#include <malloc.h>
#include <stdio.h>
#include <stdint.h>
#include <ffi.h>
float *n;
float puts_binding(char* ffs) {
	printf("hello %s\n",ffs);
	//n=(int*)malloc(sizeof(4));
	//*n=11;
	return 11.5;
}

*/
import "C"
import (
	"github.com/samber/mo"
	"runtime"
	"unsafe"
)

type NativeFunctionOption func(function *NativeFunction)
type NativeFunction struct {
	cif         C.ffi_cif
	ffiArgTypes **C.ffi_type

	Address      uintptr
	RetTypeName  RetTypeName
	ArgsTypeName []ArgTypeName
	Abi          NativeABI
	isInit       bool
}

func (n *NativeFunction) makeArgTypeNames() []*C.ffi_type {
	argTypes := n.ArgsTypeName

	if len(argTypes) == 0 {
		return nil
	}
	cargs := make([]*C.ffi_type, 0)
	for _, argType := range n.ArgsTypeName {
		cargs = append(cargs, ConvertStringTypeToFFIType(argType))
	}
	return cargs
}
func (n *NativeFunction) Call(args ...any) mo.Result[NativePointer] {
	if n.isInit == false {
		if C.ffi_prep_cif(&n.cif, n.Abi, C.uint(len(n.ArgsTypeName)), ConvertStringTypeToFFIType(n.RetTypeName), n.ffiArgTypes) != C.FFI_OK {
			return mo.Errf[NativePointer]("ffi_prep_cif失败")
		}
		if len(args) != len(n.ArgsTypeName) {
			return mo.Errf[NativePointer]("call num %d!=%d", len(args), len(n.ArgsTypeName))
		}
		n.isInit = true
	}

	var cArgValues **C.void = nil
	argValues := make([]*C.void, 0)
	argValuePins := make([]runtime.Pinner, 0)
	for i := 0; i < len(args); i++ {
		typeName := n.ArgsTypeName[i]
		argValue := ConvertAnyToFFIValue(typeName, args[i])
		argPtr := argValue
		pin := runtime.Pinner{}
		pin.Pin(argPtr)
		defer pin.Unpin()
		argValuePins = append(argValuePins, pin)
		argValues = append(argValues, (*C.void)(unsafe.Pointer(argPtr)))
	}
	if len(argValues) != 0 {
		cArgValues = &argValues[0]
	}
	argPin := runtime.Pinner{}
	argPin.Pin(cArgValues)
	defer argPin.Unpin()

	var rValue *C.void = nil
	rvaluePin := runtime.Pinner{}
	rvaluePin.Pin(&rValue)
	defer rvaluePin.Unpin()
	C.ffi_call(&n.cif, (*[0]byte)(unsafe.Pointer(n.Address)), unsafe.Pointer(&rValue), (*unsafe.Pointer)(unsafe.Pointer(cArgValues)))
	//convertRetValue := ConvertFFIRetValue(n.RetTypeName, unsafe.Pointer(&rValue))
	return mo.Ok(Ptr(unsafe.Pointer(&rValue)))
}
func NativeFunctionWithAbi(abi NativeABI) NativeFunctionOption {
	return func(function *NativeFunction) {
		function.Abi = abi
	}
}
func NewNativeFunction(ptr uintptr, retType RetTypeName, types []ArgTypeName, options ...NativeFunctionOption) *NativeFunction {
	nf := &NativeFunction{
		Address:      ptr,
		RetTypeName:  retType,
		ArgsTypeName: types,
		Abi:          DefaultAbi,
	}
	for _, option := range options {
		option(nf)
	}

	ffiArgTypes := nf.makeArgTypeNames()
	var PFFI_TYPE *C.ffi_type
	nf.ffiArgTypes = (**C.ffi_type)(C.malloc(CULong(uint64(unsafe.Sizeof(PFFI_TYPE)) * uint64(len(ffiArgTypes)))))
	sliceArgsType := unsafe.Slice(nf.ffiArgTypes, len(ffiArgTypes))
	for i := 0; i < len(ffiArgTypes); i++ {
		sliceArgsType[i] = ffiArgTypes[i]
	}

	runtime.SetFinalizer(nf, func(nb *NativeFunction) {
		C.free(unsafe.Pointer(nb.ffiArgTypes))
	})
	return nf
}
