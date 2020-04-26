package handler

import (
	"fmt"
	"testing"
	"time"

	"bitbucket.org/forfd/custm-chat/webim/models"
)

// var attrNames = "name, age, gender, tel, qq, weixin, weibo, address, email, comment"
func TestSetAttrs(t *testing.T) {
	attrs := attrsV1{
		"name":       "test",
		"age":        "20",
		"gender":     "M",
		"tel":        "1435453",
		"qq":         "31213",
		"weixin":     "weixin123",
		"weibo":      "abc",
		"address":    "adresstest",
		"email":      "abc@gmail.com",
		"comment":    "hihi",
		"测试问题":       " 测试回答",
		"vipAccount": "abc890",
	}

	visitor := &models.Visitor{}
	attrsValue := make(map[string]interface{})
	attrs.setAttrs(visitor, attrsValue)

	fmt.Printf("visitor: %+v\n", *visitor)
	fmt.Printf("attrsValue: %+v\n", attrsValue)
	fmt.Printf("attrs: %+v\n", attrs)
}

func TestSetVisitorAttrs(t *testing.T) {
	vst := &models.Visitor{
		ID:               "vID",
		EntID:            "eID",
		TraceID:          "traceID",
		Name:             "testName",
		Age:              28,
		Gender:           "M",
		Avatar:           "",
		Mobile:           "1234",
		Weibo:            "1234",
		Wechat:           "1234",
		Email:            "abc@gmail.com",
		QqNum:            "1234",
		Address:          "aaaddd",
		Remark:           "hhhh",
		VisitCnt:         0,
		VisitPageCnt:     0,
		ResidenceTimeSec: 0,
		LastVisitID:      "",
		VisitedAt:        time.Time{},
		CreatedAt:        time.Time{},
		UpdatedAt:        time.Time{},
	}

	attrs := setVisitorAttrs(vst)
	for k, v := range attrs {
		fmt.Println("key: ", k, "value: ", v)
	}
}
