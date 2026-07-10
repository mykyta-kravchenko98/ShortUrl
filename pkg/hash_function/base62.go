package hashfunction

import (
	"fmt"
	"strings"
)

const base62Chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

// DecimalToBase62 converts a decimal number to its Base62 representation
func DecimalToBase62(decimalNum int64) string {
	if decimalNum == 0 {
		return "0"
	}

	var base62Builder strings.Builder
	for decimalNum > 0 {
		remainder := decimalNum % 62
		decimalNum /= 62
		base62Builder.WriteByte(base62Chars[remainder])
	}

	// Reverse the characters to get the correct Base62 representation
	base62Bytes := []byte(base62Builder.String())
	for i, j := 0, len(base62Bytes)-1; i < j; i, j = i+1, j-1 {
		base62Bytes[i], base62Bytes[j] = base62Bytes[j], base62Bytes[i]
	}

	return string(base62Bytes)
}

func Uint64ToBase62(n uint64) string {
	if n == 0 {
		return "0"
	}

	var base62Builder strings.Builder
	for n > 0 {
		remainder := n % 62
		n /= 62
		base62Builder.WriteByte(base62Chars[remainder])
	}

	base62Bytes := []byte(base62Builder.String())
	for i, j := 0, len(base62Bytes)-1; i < j; i, j = i+1, j-1 {
		base62Bytes[i], base62Bytes[j] = base62Bytes[j], base62Bytes[i]
	}

	return string(base62Bytes)
}

// Base62ToDecimal converts a Base62 string to its decimal representation.
// Returns an error for any character outside the base62 alphabet, instead
// of silently folding it into the result - unvalidated input previously
// made strings.IndexByte return -1, which got multiplied into the running
// total without complaint, corrupting the output instead of failing loudly.
func Base62ToDecimal(base62Str string) (int64, error) {
	var decimalNum int64

	for _, char := range base62Str {
		idx := strings.IndexByte(base62Chars, byte(char))
		if idx < 0 {
			return 0, fmt.Errorf("base62: invalid character %q in %q", char, base62Str)
		}
		decimalNum = decimalNum*62 + int64(idx)
	}

	return decimalNum, nil
}
