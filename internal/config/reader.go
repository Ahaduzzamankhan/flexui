package config

import (
	"bytes"
	"io"
)

// newTrimReader strips a UTF-8 byte order mark, which editors on Windows often
// add and which encoding/json rejects.
func newTrimReader(data []byte) io.Reader {
	return bytes.NewReader(bytes.TrimPrefix(data, []byte{0xEF, 0xBB, 0xBF}))
}
