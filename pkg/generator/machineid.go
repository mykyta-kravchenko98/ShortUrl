package generator

import (
	"hash/fnv"
	"os"
	"strconv"
)

// ResolveMachineID picks the Snowflake machine ID to run with. A static,
// identical value for every replica (which is what the config file alone
// gives you) causes ID collisions the moment more than one instance runs at
// once, so this exists to give each running pod a different value without
// requiring a coordination service. Order of preference:
//
//  1. MACHINE_ID env var, if set - an explicit escape hatch for when you
//     need a guaranteed, hand-assigned value.
//  2. A hash of POD_NAME (populated from the Kubernetes downward API) modulo
//     the machine ID space (32 slots, 0-31). Deterministic per running pod,
//     no coordination needed. Collisions are possible (birthday problem
//     over 32 slots) but unlikely at the replica counts this app runs at -
//     and a collision only risks two pods sharing an ID within the same
//     millisecond+sequence-overflow window, not silent data loss.
//  3. fallback - whatever the static config file has. Only correct for a
//     single running instance (local/dev).
func ResolveMachineID(fallback int64) int64 {
	if raw := os.Getenv("MACHINE_ID"); raw != "" {
		if v, err := strconv.ParseInt(raw, 10, 64); err == nil && v >= 0 && v <= maxMachineID {
			return v
		}
	}

	if podName := os.Getenv("POD_NAME"); podName != "" {
		h := fnv.New32a()
		_, _ = h.Write([]byte(podName))
		return int64(h.Sum32() % uint32(maxMachineID+1))
	}

	return fallback
}
