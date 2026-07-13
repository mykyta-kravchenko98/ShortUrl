package cache_test

import (
	"context"
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
	ctx := context.Background()

	for _, cap := range capacities {
		// Run a sub-benchmark for each capacity
		b.Run(fmt.Sprintf("Capacity_%d", cap), func(b *testing.B) {
			c := cache.InitLRUCache(cap)

			// Ensure b.N cache operations are performed
			for n := 0; n < b.N; n++ {
				key := "key" + strconv.Itoa(n)
				value := "value" + strconv.Itoa(n)
				c.Put(ctx, key, value)
				c.Get(ctx, key)
			}

			// Request a garbage collection to release unused memory
			runtime.GC()
		})
	}
}

func Test_Put_Works_Correct_And_Not_Increment_Capacity(t *testing.T) {
	ctx := context.Background()
	to := cache.InitLRUCache(1)

	to.Put(ctx, "key1", "value1")
	to.Put(ctx, "key2", "value2")

	assert.Equal(t, 1, to.Length())
}

func Test_Get_Works_Correct_And_Return_Empty_String_If_Key_Not_Exist(t *testing.T) {
	ctx := context.Background()
	to := cache.InitLRUCache(1)
	to.Put(ctx, "key1", "value1")
	to.Put(ctx, "key2", "value2")

	result1 := to.Get(ctx, "key1")
	result2 := to.Get(ctx, "key2")

	assert.Equal(t, "", result1)
	assert.Equal(t, "value2", result2)
}
