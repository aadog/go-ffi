//go:build !windows && cgo
// +build !windows,cgo

package ffi

/*
#include <string.h>
#include <stdint.h>
*/
import "C"
import "unsafe"

func CULong(n uint64) C.ulong {
	return C.ulong(n)
}
func Memcpy(dst unsafe.Pointer, src unsafe.Pointer, l int64) {
	C.memcpy(unsafe.Pointer(dst), unsafe.Pointer(src), C.ulong(l))
}
