package common

import (
	"log"
	"testing"
)

func TestRandStringBytesMask(t *testing.T) {
	s := RandStringBytesMask(64)
	log.Println(s, len(s))

	s = RandStringBytesMask(32)
	log.Println(s, len(s))
}
