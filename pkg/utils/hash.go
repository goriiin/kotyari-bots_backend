package utils

import (
	"fmt"
	"hash/fnv"
)

func HashString(s []byte) string {
	h := fnv.New64a()
	h.Write(s)
	return fmt.Sprintf("%x", h.Sum64())
}
