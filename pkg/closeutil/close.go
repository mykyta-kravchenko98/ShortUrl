// Package closeutil gives deferred Close() calls somewhere to put the
// error instead of silently discarding it (which is what `defer f.Close()`
// does, and what errcheck/golangci-lint flags).
package closeutil

import (
	"io"
	"log/slog"
)

// Close calls c.Close() and logs a failure instead of swallowing it.
// Use it in place of a bare `defer f.Close()`:
//
//	file, err := os.Open(path)
//	if err != nil {
//	    return err
//	}
//	defer closeutil.Close(file)
func Close(c io.Closer) {
	if err := c.Close(); err != nil {
		slog.Error("failed to close resource", "error", err)
	}
}

// CloseNamed is Close with a label attached to the log line, for when a
// function defers more than one Closer and needs to tell them apart.
//
//	defer closeutil.CloseNamed("config file", file)
func CloseNamed(name string, c io.Closer) {
	if err := c.Close(); err != nil {
		slog.Error("failed to close resource", "name", name, "error", err)
	}
}
