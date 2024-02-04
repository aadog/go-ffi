//go:build !posix && cgo
// +build !posix,cgo

package ffi

/*
#include <string.h>
#include <stdint.h>
*/
import "C"
import "unsafe"

func CULong(n uint64) C.ulonglong {
	return C.ulonglong(n)
}
func Memcpy(dst unsafe.Pointer, src unsafe.Pointer, l int64) {
	C.memcpy(unsafe.Pointer(dst), unsafe.Pointer(src), C.ulonglong(l))
}
