package ffi

/*
#include <stdio.h>
#include <stdlib.h>
#include <ffi.h>
extern void NativeCallbackBinding(ffi_cif *cif, void *retVal, void **args, void *userData);
*/
import "C"
import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/samber/mo"
	"reflect"
	"runtime"
	"sync"
	"unsafe"
)

var mpCallBack sync.Map

type CallBackStruct struct {
	Fn          reflect.Value
	RetTypeName string
	ArgTypeName []string
}

//export NativeCallbackBinding
func NativeCallbackBinding(cif *C.ffi_cif, retVal unsafe.Pointer, args *unsafe.Pointer, userData unsafe.Pointer) {
	id := C.GoString((*C.char)(userData))

	iCallBack, ok := mpCallBack.Load(id)
	if !ok {
		return
	}
	callbackStruct := iCallBack.(*CallBackStruct)

	fnVal := callbackStruct.Fn
	fnType := fnVal.Type()
	fnArgCount := fnType.NumIn()
	if fnArgCount != len(callbackStruct.ArgTypeName) {
		panic(errors.New(fmt.Sprintf("NativeCallbackBinding error,cif arg num %d,func arg %d", int(cif.nargs), fnArgCount)))
	}
	if callbackStruct.RetTypeName == TVoid {
		if fnType.NumOut() != 0 {
			panic(errors.New("NativeCallbackBinding error,ret type is TVoid"))
		}
	} else {
		if fnType.NumOut() != 1 {
			panic(errors.New(fmt.Sprintf("NativeCallbackBinding error,ret type is %s", callbackStruct.RetTypeName)))
		}
	}

	fnArgs := make([]reflect.Value, 0)
	if fnArgCount > 0 {
		fnArgs = lo.FlatMap(unsafe.Slice(args, fnArgCount), func(item unsafe.Pointer, index int) []reflect.Value {
			return []reflect.Value{reflect.ValueOf(ConvertFFIValueToAny(callbackStruct.ArgTypeName[index], item))}
		})
	}
	rets := fnVal.Call(fnArgs)
	if callbackStruct.RetTypeName != TVoid {
		firstRet := rets[0]
		WriteRetValue(Ptr(uintptr(retVal)), callbackStruct.RetTypeName, firstRet.Interface())
	}
}

type NativeCallbackOption func(function *NativeCallback)
type NativeCallback struct {
	cif          C.ffi_cif
	closure      unsafe.Pointer
	bound_puts   unsafe.Pointer
	RetTypeName  RetTypeName
	ArgsTypeName []ArgTypeName
	Abi          NativeABI
	ffiArgTypes  **C.ffi_type
	id           *C.char
	maked        bool
}

func (n *NativeCallback) Ptr() unsafe.Pointer {
	return unsafe.Pointer(n.closure)
}
func (n *NativeCallback) makeArgTypeNames() []*C.ffi_type {
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

func (n *NativeCallback) MakeCall() mo.Result[uintptr] {
	if n.maked == false {
		if C.ffi_prep_cif(&n.cif, n.Abi, C.uint(len(n.ArgsTypeName)), ConvertStringTypeToFFIType(n.RetTypeName), n.ffiArgTypes) != C.FFI_OK {
			return mo.Errf[uintptr]("ffi_prep_cif失败")
		}
		if C.ffi_prep_closure_loc((*C.ffi_closure)(n.closure), &n.cif, (*[0]byte)(unsafe.Pointer(C.NativeCallbackBinding)), unsafe.Pointer(n.id), n.closure) != C.FFI_OK {
			return mo.Errf[uintptr]("ffi_prep_closure_loc error")
		}
		n.maked = true
	}
	return mo.Ok(uintptr(n.closure))
}
func NativeCallbackWithAbi(abi NativeABI) NativeCallbackOption {
	return func(function *NativeCallback) {
		function.Abi = abi
	}
}
func NewNativeCallback(fn any, retType RetTypeName, types []ArgTypeName, options ...NativeCallbackOption) *NativeCallback {
	id := uuid.NewString()
	nb := &NativeCallback{
		RetTypeName:  retType,
		ArgsTypeName: types,
		Abi:          DefaultAbi,
		id:           C.CString(id),
	}
	for _, option := range options {
		option(nb)
	}

	mpCallBack.Store(id, &CallBackStruct{
		Fn:          reflect.ValueOf(fn),
		RetTypeName: retType,
		ArgTypeName: types,
	})
	nb.closure = C.ffi_closure_alloc(CULong(uint64(unsafe.Sizeof(C.ffi_closure{}))), &nb.bound_puts)
	ffiArgTypes := nb.makeArgTypeNames()
	var PFFI_TYPE *C.ffi_type
	nb.ffiArgTypes = (**C.ffi_type)(C.malloc(CULong(uint64(unsafe.Sizeof(PFFI_TYPE)) * uint64(len(ffiArgTypes)))))
	sliceArgsType := unsafe.Slice(nb.ffiArgTypes, len(ffiArgTypes))
	for i := 0; i < len(ffiArgTypes); i++ {
		sliceArgsType[i] = ffiArgTypes[i]
	}

	runtime.SetFinalizer(nb, func(nb *NativeCallback) {
		mpCallBack.Delete(C.GoString(nb.id))
		C.free(unsafe.Pointer(nb.id))
		C.free(unsafe.Pointer(nb.ffiArgTypes))
		C.ffi_closure_free(nb.closure)
	})
	return nb
}
