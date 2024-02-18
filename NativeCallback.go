package ffi

/*
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <ffi.h>
extern void goFFIBinding(ffi_cif *cif,void *ret,void* args[],void *userData);
extern void NativeCallbackBinding(ffi_cif *cif, void *retVal, void **args, void *userData);
*/
import "C"
import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"reflect"
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
		sliceArgs := unsafe.Slice(args, fnArgCount)
		for index, item := range sliceArgs {
			typeName := callbackStruct.ArgTypeName[index]
			anyValue := ConvertFFIValueToAny(typeName, item)
			fnArgs = append(fnArgs, reflect.ValueOf(anyValue))
		}
	}

	rets := fnVal.Call(fnArgs)
	if callbackStruct.RetTypeName != TVoid {
		firstRet := rets[0]
		b := firstRet.Interface()
		WriteRetValue(Ptr(unsafe.Pointer(retVal)), callbackStruct.RetTypeName, b)
	}

	return
}

type NativeCallbackOption func(function *NativeCallback)
type NativeCallback struct {
	cif          *C.ffi_cif
	closure      *C.ffi_closure
	bound_puts   unsafe.Pointer
	RetTypeName  RetTypeName
	ArgsTypeName []ArgTypeName
	Abi          NativeABI
	id           unsafe.Pointer
	maked        bool
}

func (n *NativeCallback) Closure() unsafe.Pointer {
	return unsafe.Pointer(n.bound_puts)
}
func (n *NativeCallback) Ptr() unsafe.Pointer {
	return unsafe.Pointer(n.bound_puts)
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

func (n *NativeCallback) Free() {
	
	C.ffi_closure_free(unsafe.Pointer((unsafe.Pointer(n.closure))))
	fmt.Println("free")

	C.free(unsafe.Pointer(n.cif.arg_types))
	C.free(unsafe.Pointer(n.cif))
	C.free(n.id)
	mpCallBack.Delete(C.GoString((*C.char)(n.id)))
}

func NativeCallbackWithAbi(abi NativeABI) NativeCallbackOption {
	return func(function *NativeCallback) {
		function.Abi = abi
	}
}
func NewNativeCallback(fn any, retType RetTypeName, types []ArgTypeName, options ...NativeCallbackOption) *NativeCallback {
	nb := &NativeCallback{
		RetTypeName:  retType,
		ArgsTypeName: types,
		Abi:          DefaultAbi,
		id:           unsafe.Pointer(C.CString(uuid.New().String())),
	}

	for _, option := range options {
		option(nb)
	}

	mpCallBack.Store(C.GoString((*C.char)(nb.id)), &CallBackStruct{
		Fn:          reflect.ValueOf(fn),
		RetTypeName: retType,
		ArgTypeName: types,
	})
	nb.cif = (*C.ffi_cif)(C.malloc(C.size_t(unsafe.Sizeof(C.ffi_cif{}))))
	c := C.ffi_closure{}

	nb.closure = (*C.ffi_closure)(C.ffi_closure_alloc(C.size_t(unsafe.Sizeof(c)), &nb.bound_puts))
	//
	tp := C.ffi_type{}
	ffiArgTypes := (**C.ffi_type)(C.malloc(C.size_t(unsafe.Sizeof(unsafe.Pointer(&tp))) * C.size_t(len(types))))
	sliceArgsType := unsafe.Slice(ffiArgTypes, len(types))
	for i := 0; i < len(types); i++ {
		sliceArgsType[i] = ConvertStringTypeToFFIType(types[i])
	}

	if status := Status(C.ffi_prep_cif(nb.cif, C.FFI_DEFAULT_ABI, C.uint(len(nb.ArgsTypeName)), ConvertStringTypeToFFIType(nb.RetTypeName), ffiArgTypes)); status != OK {
		panic(errors.New(fmt.Sprintf("%d", status)))
	}

	if C.ffi_prep_closure_loc((*C.ffi_closure)(nb.closure), nb.cif, (*[0]byte)(unsafe.Pointer(C.NativeCallbackBinding)), nb.id, unsafe.Pointer(nb.bound_puts)) != C.FFI_OK {
		panic(errors.New("ffi_prep_closure_loc error"))
	}
	return nb
}
