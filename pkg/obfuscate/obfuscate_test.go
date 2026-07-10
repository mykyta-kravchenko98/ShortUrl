package obfuscate_test

import (
	"testing"

	"github.com/mykyta-kravchenko98/ShortUrl/pkg/obfuscate"
	"github.com/stretchr/testify/assert"
)

const testKey uint64 = 0x1234567890ABCDEF

func Test_Mix_Is_Deterministic(t *testing.T) {
	var id int64 = 7087840649386528768
	assert.Equal(t, obfuscate.Mix(id, testKey), obfuscate.Mix(id, testKey))
}

func Test_Mix_Is_Keyed(t *testing.T) {
	var id int64 = 7087840649386528768
	assert.NotEqual(t, obfuscate.Mix(id, testKey), obfuscate.Mix(id, testKey+1))
}

func Test_Mix_Sequential_Inputs_Do_Not_Collide(t *testing.T) {
	const base int64 = 7087840649386528000
	seen := make(map[uint64]int64, 4096)

	for i := int64(0); i < 4096; i++ {
		out := obfuscate.Mix(base+i, testKey)
		if prev, ok := seen[out]; ok {
			t.Fatalf("collision: id %d and id %d both mixed to %d", prev, base+i, out)
		}
		seen[out] = base + i
	}
}

func Test_Mix_Sequential_Inputs_Avalanche(t *testing.T) {
	a := obfuscate.Mix(1000, testKey)
	b := obfuscate.Mix(1001, testKey)

	diffBits := popcount(a ^ b)
	// A single-bit input change should flip roughly half of 64 output
	// bits under a good mix. Assert it's at least "a lot", not "off by a
	// bit or two" (which is what a plain XOR-only mix would give).
	assert.Greater(t, diffBits, 16, "expected substantial bit diffusion between mixed(1000) and mixed(1001)")
}

func popcount(x uint64) int {
	count := 0
	for x != 0 {
		count += int(x & 1)
		x >>= 1
	}
	return count
}
