package ffi

import (
	"fmt"
	"testing"
)

func TestSingle(t *testing.T) {
	nc := NewNativeCallback(func(a int) int {
		fmt.Println("param:", a)
		return 222
	}, Tint, []ArgTypeName{Tint})
	nf := NewNativeFunction(nc.MakeCall().MustGet(), Tint, []ArgTypeName{Tint})
	fmt.Printf("return:%d\n", nf.Call(1).MustGet().ReadInt())
}
func TestCall(t *testing.T) {
	for i := 0; i < 100000; i++ {
		fn := NewNativeCallback(func(a int) int {
			//fmt.Println(a)
			return 222
		}, Tint, []ArgTypeName{Tint})
		fn1 := NewNativeFunction(fn.MakeCall().MustGet(), Tint, []ArgTypeName{Tint})
		fmt.Printf("return:%d\n", fn1.Call(1).MustGet().ReadInt())
	}
	fmt.Println("完成")
}
