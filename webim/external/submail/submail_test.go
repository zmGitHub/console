package submail

import (
	"testing"

	"bitbucket.org/forfd/custm-chat/webim/common"
	"bitbucket.org/forfd/custm-chat/webim/common/logging"
	"bitbucket.org/forfd/custm-chat/webim/conf"
)

func TestSendWithTmpl(t *testing.T) {
	log.NewLogging()
	NewClient(&conf.SubMailConfig{AppID: "14028", AppKey: "33d307f1723821ba9ccad64e6894a628", SignType: "sha1"})

	link := common.EncodeURL("http://backendtest", map[string]string{
		"email": "someone@qq.com",
		"name":  "someone",
	})
	err := sendEmail("2550418657@qq.com", "test", "TEST", link)
	log.Logger.Println(err)
}
