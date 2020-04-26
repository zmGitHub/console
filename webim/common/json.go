package common

import (
	"bytes"

	"github.com/json-iterator/go"
)

var json = jsoniter.Config{
	EscapeHTML:             false,
	SortMapKeys:            true,
	ValidateJsonRawMessage: true,
}.Froze()

func Marshal(object interface{}) (string, error) {
	return json.MarshalToString(object)
}

func Unmarshal(s string, i interface{}) error {
	return json.UnmarshalFromString(s, i)
}

func MarshalUnescape(object interface{}) (string, error) {
	bs, err := json.Marshal(object)
	if err != nil {
		return "", err
	}

	bs = ReplaceHtmlChar(bs)
	return string(bs), nil
}

func ReplaceHtmlChar(data []byte) []byte {
	data = bytes.Replace(data, []byte("\\u0026"), []byte("&"), -1)
	data = bytes.Replace(data, []byte("\\u003c"), []byte("<"), -1)
	data = bytes.Replace(data, []byte("\\u003e"), []byte(">"), -1)
	data = bytes.Replace(data, []byte("\\u003d"), []byte("="), -1)
	return data
}
