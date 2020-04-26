package handler

import (
	"log"
	"net/url"
	"testing"
)

func TestIMService_QrCode(t *testing.T) {
	text := `http%3A%2F%2Fstatic.chat186.com%2Fdist%2Fstandalone.html%3Feid%3Dbjstn2pfua60sjh9jae0`
	link, err := url.QueryUnescape(text)
	log.Println(link, err)
}
