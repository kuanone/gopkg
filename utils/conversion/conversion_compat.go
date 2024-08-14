//go:build !go1.15 && !go1.16 && !go1.17 && !go1.18 && !go1.19 && !go1.20 && !go1.21 && !go1.22 && !go1.23
// +build !go1.15,!go1.16,!go1.17,!go1.18,!go1.19,!go1.20,!go1.21,!go1.22,!go1.23

package conversion

import "unsafe"

func toBytesImpl(s string) []byte {
	if len(s) == 0 {
		return nil
	}

	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

func toStringImpl(b []byte) string {
	if len(b) == 0 {
		return ""
	}
	return *(*string)(unsafe.Pointer(&b))
}
