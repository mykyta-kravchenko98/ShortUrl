package hashfunction_test

import (
	"testing"

	hashfunction "github.com/mykyta-kravchenko98/ShortUrl/pkg/hash_function"
	"github.com/stretchr/testify/assert"
)

const (
	calculatedId   int64 = 7087840649386528768
	calculatedHash       = "8RaNtMde0cS"
)

func Test_Base62_Encode_In_Decimal_Is_Correct(t *testing.T) {
	result := hashfunction.Base62ToDecimal(calculatedHash)

	assert.Equal(t, result, calculatedId)
}

func Test_Base62_Decode_In_Hash_Is_Correct(t *testing.T) {
	result := hashfunction.DecimalToBase62(calculatedId)

	assert.Equal(t, result, calculatedHash)
}

func Benchmark_Base62_Decode(b *testing.B) {
	// Ensure b.N cache operations are performed
	for n := 0; n < b.N; n++ {
		hashfunction.DecimalToBase62(calculatedId)
	}
}

func Benchmark_Base62_Encode(b *testing.B) {
	// Ensure b.N cache operations are performed
	for n := 0; n < b.N; n++ {
		hashfunction.Base62ToDecimal(calculatedHash)
	}
}
