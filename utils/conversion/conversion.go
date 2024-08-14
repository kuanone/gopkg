package conversion

// ToBytes converts a string to a byte slice.
func ToBytes(s string) []byte {
	if len(s) == 0 {
		return nil
	}
	return toBytesImpl(s)
}

// ToString converts a byte slice to a string.
func ToString(b []byte) string {
	if len(b) == 0 {
		return ""
	}
	return toStringImpl(b)
}
