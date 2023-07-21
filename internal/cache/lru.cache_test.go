package cache_test

import (
	"fmt"
	"runtime"
	"strconv"
	"testing"

	"github.com/mykyta-kravchenko98/ShortUrl/internal/cache"
	"github.com/stretchr/testify/assert"
)

func BenchmarkLRUCacheMemory(b *testing.B) {
	// Define different cache capacities to benchmark
	capacities := []int{1000, 5000, 10000, 50000, 100000, 250000, 500000, 1000000, 2500000, 10000000}

	for _, cap := range capacities {
		// Run a sub-benchmark for each capacity
		b.Run(fmt.Sprintf("Capacity_%d", cap), func(b *testing.B) {
			c := cache.InitLRUCache(cap)

			// Ensure b.N cache operations are performed
			for n := 0; n < b.N; n++ {
				key := "key" + strconv.Itoa(n)
				value := "value" + strconv.Itoa(n)
				c.Put(key, value)
				c.Get(key)
			}

			// Request a garbage collection to release unused memory
			runtime.GC()
		})
	}
}

func Test_Put_Works_Correct_And_Not_Increment_Capacity(t *testing.T) {
	to := cache.InitLRUCache(1)

	to.Put("key1", "value1")
	to.Put("key2", "value2")

	assert.Equal(t, 1, to.Length())
}

func Test_Get_Works_Correct_And_Return_Empty_String_If_Key_Not_Exist(t *testing.T) {
	to := cache.InitLRUCache(1)
	to.Put("key1", "value1")
	to.Put("key2", "value2")

	result1 := to.Get("key1")
	result2 := to.Get("key2")

	assert.Equal(t, "", result1)
	assert.Equal(t, "value2", result2)
}
