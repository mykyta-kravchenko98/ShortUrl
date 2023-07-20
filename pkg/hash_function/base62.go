package hashfunction

import "strings"

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

// Base62ToDecimal converts a Base62 string to its decimal representation
func Base62ToDecimal(base62Str string) int64 {
	var decimalNum int64

	for _, char := range base62Str {
		decimalNum = decimalNum*62 + int64(strings.IndexByte(base62Chars, byte(char)))
	}

	return decimalNum
}
