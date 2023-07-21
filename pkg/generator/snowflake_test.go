package generator

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_Normal_Behavior(t *testing.T) {
	to := build_NewSnowflake()

	result, err := to.NextID()
	assert.NoError(t, err)
	assert.NotEqual(t, int64(0), result)
}

// Adjust the current timestamp to be a smaller value than the previous one
// The next call to NextID should return an error
func Test_Backward_Clock_Movement(t *testing.T) {
	to := build_NewSnowflake()

	to.lastTimestamp = time.Now().Add(1*time.Minute).UnixNano() / int64(time.Millisecond)
	result, err := to.NextID()

	assert.Error(t, err)
	assert.Equal(t, int64(0), result)
}

// Set the sequence number to the maximum value (maxSequence) and the timestamp to the current time
// Subsequent calls to NextID should wait for the next millisecond to avoid sequence number overflow
func Test_Sequence_Number_Overflow(t *testing.T) {
	to := build_NewSnowflake()

	to.sequence = maxSequence
	time.Sleep(1 * time.Millisecond) // Wait for the next millisecond
	id3, err := to.NextID()

	assert.NoError(t, err)
	assert.NotEqual(t, int64(0), id3)
}

func Test_Snowflage_DatacenterId_Can_Not_Be_Great_Then_31(t *testing.T) {
	_, err := NewSnowflake(32, 1)

	assert.Error(t, err)
}

func Test_Snowflage_DatacenterId_Can_Not_Be_Less_Then_0(t *testing.T) {
	_, err := NewSnowflake(-1, 1)

	assert.Error(t, err)
}

func Test_Snowflage_MachineId_Can_Not_Be_Great_Then_31(t *testing.T) {
	_, err := NewSnowflake(1, 32)

	assert.Error(t, err)
}

func Test_Snowflage_MachineId_Can_Not_Be_Less_Then_0(t *testing.T) {
	_, err := NewSnowflake(1, -1)

	assert.Error(t, err)
}

func Benchmark_Snowflake_Generation(b *testing.B) {
	to, _ := NewSnowflake(0, 0)
	// Ensure b.N cache operations are performed
	for n := 0; n < b.N; n++ {
		_, err := to.NextID()
		if err != nil {
			b.Fatalf("Error generating ID: %v", err)
		}
	}
}

func build_NewSnowflake() *snowflake {
	return &snowflake{
		lastTimestamp: 0,
		sequence:      0,
		datacenterID:  1,
		machineID:     1,
	}
}
