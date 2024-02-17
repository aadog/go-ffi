#include "goffi.h"
extern void NativeCallbackBinding(ffi_cif *cif, void *retVal, void **args, void *userData);
void goFFIBinding(ffi_cif *cif,void *ret,void* args[],void *userData){
    *(ffi_arg*)(ret)=1;
    return;
    NativeCallbackBinding(cif,ret,args,userData);
}
