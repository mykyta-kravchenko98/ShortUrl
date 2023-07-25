package generator

import (
	"errors"
	"sync"
	"time"
)

const (
	// The number of bits allocated for each component in the Snowflake ID
	timestampBits    = 41
	datacenterIDBits = 5
	machineIDBits    = 5
	sequenceBits     = 12

	// The maximum values for datacenter ID, machine ID, and sequence number
	maxDatacenterID = -1 ^ (-1 << datacenterIDBits)
	maxMachineID    = -1 ^ (-1 << machineIDBits)
	maxSequence     = -1 ^ (-1 << sequenceBits)

	// The shift values for each component
	timestampShift    = datacenterIDBits + machineIDBits + sequenceBits
	datacenterIDShift = machineIDBits + sequenceBits
	machineIDShift    = sequenceBits
)

// snowflake struct to hold the state of the Snowflake generator
type snowflake struct {
	datacenterID  int64
	machineID     int64
	lastTimestamp int64
	sequence      int64
	lock          sync.Mutex
}

// Snowflake - interface for interaction with snowflake generator implementation
type Snowflake interface {
	NextID() (int64, error)
}

// NewSnowflake creates a new snowflake instance with the given datacenter ID and machine ID and return Snowflake interface
func NewSnowflake(datacenterID, machineID int64) (Snowflake, error) {
	if datacenterID < 0 || datacenterID > maxDatacenterID || machineID < 0 || machineID > maxMachineID {
		return nil, errors.New("datacenter ID or machine ID out of range")
	}
	return &snowflake{
		datacenterID: datacenterID,
		machineID:    machineID,
	}, nil
}

// NextID generates the next unique ID using the Snowflake algorithm
func (s *snowflake) NextID() (int64, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	timestamp := time.Now().UnixNano() / int64(time.Millisecond)

	if timestamp < s.lastTimestamp {
		return 0, errors.New("clock moved backwards")
	}

	if timestamp == s.lastTimestamp {
		s.sequence = (s.sequence + 1) & maxSequence
		if s.sequence == 0 {
			// Sequence number overflows, wait until the next millisecond
			for timestamp == s.lastTimestamp {
				timestamp = time.Now().UnixNano() / int64(time.Millisecond)
			}
		}
	} else {
		s.sequence = 0
	}

	s.lastTimestamp = timestamp

	id := ((timestamp << timestampShift) | (s.datacenterID << datacenterIDShift) | (s.machineID << machineIDShift) | s.sequence)
	return id, nil
}
