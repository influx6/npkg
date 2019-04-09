package nunsafe

import "unsafe"

//*****************************************************
// unsafe methods
//*****************************************************

// Bytes2String uses a unsafe.Pointer trick to convert a byteslice into
// a string, becare not to modify the slice or string, has it will
// reflect on either.
func Bytes2String(bc []byte) string {
	return *(*string)(unsafe.Pointer(&bc))
}

// String2Bytes uses a unsafe.Pointer trick to convert a string into
// a byte slice, becare not to modify the slice or string, has
// it will reflect on either.
func String2Bytes(bc string) []byte {
	return *(*[]byte)(unsafe.Pointer(&bc))
}

// Noescape hides a pointer from escape analysis.  noescape is
// the identity function but escape analysis doesn't think the
// output depends on the input. noescape is inlined and currently
// compiles down to zero instructions.
// USE CAREFULLY!
// This was copied from the runtime; see issues 23382 and 7921.
//go:nosplit
func Noescape(p unsafe.Pointer) unsafe.Pointer {
	x := uintptr(p)
	return unsafe.Pointer(x ^ 0)
}
