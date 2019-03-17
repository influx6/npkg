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
