package ffi

import "C"
import "unsafe"

func GoString(p *C.char) string {
	return C.GoString(p)
}
func copyBytes(src unsafe.Pointer, strLen int) []byte {
	if strLen == 0 {
		return nil
	}
	str := make([]byte, strLen)
	for i := 0; i < strLen; i++ {
		str[i] = *(*byte)(unsafe.Pointer(uintptr(src) + uintptr(i)))
	}
	return str
}
func writeBytes(src unsafe.Pointer, bt []byte) {
	if len(bt) == 0 {
		return
	}
	for i := 0; i < len(bt); i++ {
		*(*byte)(unsafe.Pointer(uintptr(src) + uintptr(i))) = bt[i]
	}
}
