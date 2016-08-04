package utils

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
)

func Round(v float32) int {
	if v < 0 {
		return int(v - 0.5)
	} else {
		return int(v + 0.4999999)
	}
}

func Sha1Sign(data string) string {
	t := sha1.New()
	io.WriteString(t, data)
	return hex.EncodeToString(t.Sum(nil))
}
