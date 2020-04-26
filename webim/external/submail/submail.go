package submail

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/parnurzeal/gorequest"

	log "bitbucket.org/forfd/custm-chat/webim/common/logging"
	"bitbucket.org/forfd/custm-chat/webim/conf"
)

var (
	sendMailEndpoint = `https://api.mysubmail.com/mail/send.json`
	xSendMail        = `https://api.mysubmail.com/mail/xsend.json`
	from             = `no-reply@chat186.com`
	SubmailClient    *Client
)

type Client struct {
	appID    string
	appKey   string
	signType string
}

func NewClient(c *conf.SubMailConfig) *Client {
	cli := &Client{
		appID:    c.AppID,
		appKey:   c.AppKey,
		signType: c.SignType,
	}
	SubmailClient = cli
	return cli
}

func (c *Client) SendEmail(to, subject, content string) error {
	name := strings.Split(to, "@")[0]
	request := gorequest.New()
	// appid=your_app_id
	//&to=leo <leo@submail.cn>,retro@submail.cn
	//&subject=testing_Subject
	//&text=testing_text_body
	//&from=no-reply@submail.cn
	//&signature=your_app_key
	requestContent := `appid=%s&to=%s <%s>&subject=%s&text=%s&from=%s&signature=%s`
	requestContent = fmt.Sprintf(requestContent, c.appID, name, to, subject, content, from, c.appKey)
	_, body, errs := request.Post(sendMailEndpoint).
		Send(requestContent).
		End()

	log.Logger.Infof("send mail result: %s", body)
	if len(errs) > 0 {
		return errs[0]
	}

	return nil
}

func (c *Client) SendMailWithHtml(to, subject, text, htmlContent string) error {
	name := strings.Split(to, "@")[0]
	request := gorequest.New()
	// appid=your_app_id
	//&to=leo <leo@submail.cn>,retro@submail.cn
	//&subject=testing_Subject
	//&text=testing_text_body
	//&from=no-reply@submail.cn
	//&signature=your_app_key
	requestContent := `appid=%s&to=%s <%s>&subject=%s&text=%s&html=%s&from=%s&signature=%s`
	requestContent = fmt.Sprintf(requestContent, c.appID, name, to, subject, text, htmlContent, from, c.appKey)
	_, body, errs := request.Post(sendMailEndpoint).
		Send(requestContent).
		End()

	log.Logger.Infof("send mail result: %s", body)
	if len(errs) > 0 {
		return errs[0]
	}

	return nil
}

func SendActivateEmail(to, subject, text, link string) error {
	return SendEmail(to, subject, text, link)
}

// MFic23
func SendActivateEnterprise(to, subject, text, link string) error {
	return SendEmail(to, subject, text, link)
}

func SendEmail(to, subject, text, link string) error {
	name := strings.Split(to, "@")[0]
	request := gorequest.New()

	tmpl := `appid=%s&to=%s <%s>&subject=%s&text=%s&from=%s&signature=%s`
	encodeLink := url.QueryEscape(link)
	log.Logger.Infof("encodeLink:%s", encodeLink)

	content := fmt.Sprintf(`
%s
%s
`, text, encodeLink)
	data := fmt.Sprintf(
		tmpl, SubmailClient.appID, name, to, subject, content, from, SubmailClient.appKey,
	)
	log.Logger.Infof("post content: %s", data)

	_, body, errs := request.Post(sendMailEndpoint).
		Set("Content-Type", "application/x-www-form-urlencoded").
		Send(data).
		End()

	log.Logger.Infof("send mail result: %s", body)
	if len(errs) > 0 {
		return errs[0]
	}

	return nil
}
