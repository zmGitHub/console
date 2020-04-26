package test

import "log"

type MailClient struct {
	Content string
}

func (cli *MailClient) SendEmail(address, title, content string) error {
	log.Println(address)
	log.Println(title)
	log.Println(content)
	cli.Content = content
	return nil
}

func (cli *MailClient) SendMailWithHtml(to, subject, text, htmlContent string) error {
	return nil
}
