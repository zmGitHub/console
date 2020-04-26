package common

import (
	"log"
	"testing"
)

func TestReplaceHtmlChar(t *testing.T) {
	data := `\u003cdiv\u003e\u003cspan style=\"color: rgb(87, 166, 255)\"\u003e\u003cem\u003eNew\u0026nbsp;peomotions\u0026nbsp;msgs\u0026nbsp;42432432\u003c/em\u003e\u003c/span\u003e\u003c/div\u003e`
	bs := ReplaceHtmlChar([]byte(data))
	log.Println(string(bs))
}
