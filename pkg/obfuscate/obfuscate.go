package obfuscate

func Mix(id int64, key uint64) uint64 {
	x := uint64(id) ^ key
	x += 0x9E3779B97F4A7C15
	x = (x ^ (x >> 30)) * 0xBF58476D1CE4E5B9
	x = (x ^ (x >> 27)) * 0x94D049BB133111EB
	x ^= x >> 31
	return x
}
