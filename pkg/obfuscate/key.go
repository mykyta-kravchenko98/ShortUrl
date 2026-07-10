package obfuscate

import (
	"crypto/rand"
	"encoding/binary"
	"log/slog"
	"os"
	"strconv"
)

func KeyFromEnv(envVar string) uint64 {
	if v := os.Getenv(envVar); v != "" {
		key, err := strconv.ParseUint(v, 16, 64)
		if err == nil {
			return key
		}
		slog.Error("invalid value for env var, expected 16 hex chars - falling back to a random key",
			"envVar", envVar, "error", err)
	}

	var buf [8]byte
	if _, err := rand.Read(buf[:]); err != nil {
		// crypto/rand failing means the OS entropy source is broken; there
		// is no sane fallback that's still "unpredictable".
		panic("obfuscate: failed to generate random key: " + err.Error())
	}

	slog.Warn("id obfuscation key not set, generated an ephemeral one for this process - short_url generation will look different after every restart, which is expected and fine",
		"envVar", envVar)

	return binary.BigEndian.Uint64(buf[:])
}
