package handler

import (
	"log"
	"testing"

	"bitbucket.org/forfd/custm-chat/webim/test"
)

func TestClearTables(t *testing.T) {
	test.InitTest()
	defer test.Clear()
	log.Println("clear tables...")
}
