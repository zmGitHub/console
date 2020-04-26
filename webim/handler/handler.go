package handler

import (
	"io"
	"time"

	"bitbucket.org/forfd/custm-chat/webim/auth"
	"bitbucket.org/forfd/custm-chat/webim/external/timedevent"
	"bitbucket.org/forfd/custm-chat/webim/imclient"
)

type EMailSender interface {
	SendEmail(address, title, content string) error
	SendMailWithHtml(to, subject, text, htmlContent string) error
}

type Uploader interface {
	Upload(fileName string, reader io.Reader, expireAt time.Time) (url string, err error)
}

// IPGeoLocation get location info of ip address
type IPGeoLocation interface {
	GetLocation(ip string) (country, province, city, isp string, err error)
}

type TaskHandler interface {
	AddSendingEntMessage(req *timedevent.AddSendingEntMessageReq) error
	AddEndingConversation(req *timedevent.AddEndingConversationTaskReq) error
	AddSendingNoRespMessage(req *timedevent.AddSendingNoRespMessageReq) error
	AddOfflineEndConversation(req *timedevent.AddOfflineTaskReq) error
	DeleteJob(req *timedevent.DeleteTaskReq) error
}

// IMService implements all apis of on line chat service
type IMService struct {
	mailClient  EMailSender
	imCli       imclient.IMClient
	auth        auth.MultiAgentsAuthenticate
	uploader    Uploader
	loc         IPGeoLocation
	taskHandler TaskHandler
}

// Option ImService option
type Option func(s *IMService)

// WithIPGeoLocation set IPGeoLocation field
func WithIPGeoLocation(loc IPGeoLocation) Option {
	return func(s *IMService) {
		s.loc = loc
	}
}

// WithEmailSender set mail client
func WithEmailSender(sender EMailSender) Option {
	return func(s *IMService) {
		s.mailClient = sender
	}
}

func WithIMClient(cli imclient.IMClient) Option {
	return func(s *IMService) {
		s.imCli = cli
	}
}

func WithAuth(auther auth.MultiAgentsAuthenticate) Option {
	return func(s *IMService) {
		s.auth = auther
	}
}

func WithUploader(uploader Uploader) Option {
	return func(s *IMService) {
		s.uploader = uploader
	}
}

func WithTaskHandler(handler TaskHandler) Option {
	return func(s *IMService) {
		s.taskHandler = handler
	}
}

// NewIMService create handler with options
func NewIMService(opts ...Option) *IMService {
	s := &IMService{}
	for _, opt := range opts {
		opt(s)
	}

	return s
}
