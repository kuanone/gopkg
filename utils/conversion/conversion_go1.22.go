//go:build go1.22 && !go1.23
// +build go1.22,!go1.23

package conversion

import "unsafe"

func toBytesImpl(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

func toStringImpl(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
}
