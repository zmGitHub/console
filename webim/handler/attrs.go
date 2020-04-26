package handler

import (
	"strconv"
	"strings"

	"bitbucket.org/forfd/custm-chat/webim/models"
)

var attrNames = "name, age, gender, tel, qq, weixin, weibo, address, email, comment"

type attrsV1 map[string]string

// Name:    v.Name,
//		Age:     v.Age,
//		Gender:  v.Gender,
//		Tel:     v.Mobile,
//		QQ:      v.QqNum,
//		Weixin:  v.Wechat,
//		Weibo:   v.Weibo,
//		Address: v.Address,
//		Email:   v.Email,
//		Comment: v.Remark,
func setVisitorAttrs(visitor *models.Visitor) map[string]interface{} {
	attrs := make(map[string]interface{}, 10)

	attrs["name"] = visitor.Name
	attrs["age"] = visitor.Age
	attrs["gender"] = visitor.Gender
	attrs["tel"] = visitor.Mobile
	attrs["qq"] = visitor.QqNum
	attrs["weixin"] = visitor.Wechat
	attrs["weibo"] = visitor.Weibo
	attrs["address"] = visitor.Address
	attrs["email"] = visitor.Email
	attrs["comment"] = visitor.Remark

	return attrs
}

// Name    string `json:"name"`
//	Age     string `json:"age"`
//	Gender  string `json:"gender"`
//	Tel     string `json:"tel"`    // tel
//	QQ      string `json:"qq"`     // qq
//	Weixin  string `json:"weixin"` // weixin, weibo, address, email, comment
//	Weibo   string `json:"weibo"`
//	Address string `json:"address"`
//	Email   string `json:"email"`   // email
//	Comment string `json:"comment"` // comment
func (a attrsV1) getValue(key string) string {
	value := a[key]
	delete(a, key)
	return value
}

func (a attrsV1) setAge(visitor *models.Visitor, attrsValue map[string]interface{}) {
	age := a.getValue("age")
	if age != "" {
		ageInt, err := strconv.Atoi(age)
		if err == nil {
			visitor.Age = ageInt
			attrsValue["age"] = age
		}
	}
}

func (a attrsV1) setName(visitor *models.Visitor, attrsValue map[string]interface{}) {
	name := a.getValue("name")
	if name != "" {
		visitor.Name = name
		attrsValue["name"] = name
	}
}

func (a attrsV1) setGender(visitor *models.Visitor, attrsValue map[string]interface{}) {
	gender := a.getValue("gender")
	if gender != "" {
		visitor.Gender = gender
		attrsValue["gender"] = gender
	}
}

func (a attrsV1) setTel(visitor *models.Visitor, attrsValue map[string]interface{}) {
	tel := a.getValue("tel")
	if tel != "" {
		visitor.Mobile = tel
		attrsValue["tel"] = tel
	}
}

func (a attrsV1) setAttrs(visitor *models.Visitor, attrsValue map[string]interface{}) {
	attrNames := "age, name, gender, tel, qq, weixin, weibo, address, email, comment"
	fields := strings.Split(attrNames, ", ")
	for _, field := range fields {
		switch field {
		case "age":
			a.setAge(visitor, attrsValue)
		case "name":
			a.setName(visitor, attrsValue)
		case "gender":
			a.setGender(visitor, attrsValue)
		case "tel":
			a.setTel(visitor, attrsValue)
		case "qq":
			qq := a.getValue("qq")
			if qq != "" {
				visitor.QqNum = qq
				attrsValue["qq"] = qq
			}
		case "weixin":
			wechat := a.getValue("weixin")
			if wechat != "" {
				visitor.Wechat = wechat
				attrsValue["weixin"] = wechat
			}
		case "weibo":
			weibo := a.getValue("weibo")
			if weibo != "" {
				visitor.Weibo = weibo
				attrsValue["weibo"] = weibo
			}
		case "address":
			address := a.getValue("address")
			if address != "" {
				visitor.Address = address
				attrsValue["address"] = address
			}
		case "email":
			email := a.getValue("email")
			if email != "" {
				visitor.Email = email
				attrsValue["email"] = email
			}
		case "comment":
			remark := a.getValue("comment")
			if remark != "" {
				visitor.Remark = remark
				attrsValue["comment"] = remark
			}
		}
	}
}
