package common

import (
	"fmt"
	"log"
	"testing"
)

func TestGetHostFromURL(t *testing.T) {
	url := "http://3.0.248.6:8090/signin"
	log.Println(GetHostFromURL(url))
}

func TestEncodeURL(t *testing.T) {
	base := `http://localhost:8080`
	values := map[string]string{
		"token":    "tttt123tttt",
		"end_id":   "abcd123",
		"email":    "abcde@qq.com",
		"fullname": "测试企业",
	}
	url := EncodeURL(base, values)
	fmt.Println(url)
}
