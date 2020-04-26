package sendcloud

import (
	"bytes"
	"net/http"
	"net/url"

	"bitbucket.org/forfd/custm-chat/webim/common/logging"
	"bitbucket.org/forfd/custm-chat/webim/conf"
)

type Client struct {
	endPoint string
	apiUser  string
	apiKey   string
	from     string
	fromName string
}

// params = {
//    "apiUser": API_USER, # 使用api_user和api_key进行验证
//    "apiKey" : API_KEY,
//    "to" : "test@ifaxin.com", # 收件人地址, 用正确邮件地址替代, 多个地址用';'分隔
//    "from" : "sendcloud@sendcloud.org", # 发信人, 用正确邮件地址替代
//    "fromName" : "SendCloud",
//    "subject" : "SendCloud python common",
//    "html": "欢迎使用SendCloud"
//}

func NewClient(config *conf.EmailConfig) *Client {
	return &Client{
		endPoint: config.EndPoint,
		apiUser:  config.APIUser,
		apiKey:   config.APIKey,
		from:     config.From,
		fromName: config.From,
	}
}

// http://www.sendcloud.net/doc/email_v2/code/#goa
func (cli *Client) SendEmail(to, subject, content string) error {
	RequestURI := cli.endPoint
	PostParams := url.Values{
		"apiUser":  {cli.apiUser},
		"apiKey":   {cli.apiKey},
		"from":     {cli.from},
		"fromName": {cli.fromName},
		"to":       {to},
		"subject":  {subject},
		"html":     {content},
	}
	PostBody := bytes.NewBufferString(PostParams.Encode())
	_, err := http.Post(RequestURI, "application/x-www-form-urlencoded", PostBody)
	if err != nil {
		log.Logger.Errorf("send mail error: %v", err)
		return err
	}

	return nil
}
