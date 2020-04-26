package common

import (
	"log"
	"testing"
)

func TestGenUniqueID(t *testing.T) {
	id := GenUniqueID()
	log.Println(id, len(id))
}
