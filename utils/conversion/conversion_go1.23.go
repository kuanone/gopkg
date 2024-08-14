//go:build go1.23
// +build go1.23

package conversion

func toBytesImpl(s string) []byte {
	return []byte(s)
}

func toStringImpl(b []byte) string {
	return string(b)
}
