package imclient

import (
	"context"
	"fmt"
	"time"

	"github.com/avast/retry-go"
	"github.com/centrifugal/gocent"

	"bitbucket.org/forfd/custm-chat/webim/conf"
)

var CentriClient *CentrifugoClient

type CentrifugoClient struct {
	cli              *gocent.Client
	onlineAgentsChan string
	retryTimes       uint
}

func InitClient(conf *conf.CentrifugoConfig) {
	CentriClient = NewClient(conf)
}

func NewClient(conf *conf.CentrifugoConfig) *CentrifugoClient {
	return &CentrifugoClient{
		cli: gocent.New(gocent.Config{
			Addr: conf.API,
			Key:  conf.AuthKey,
		}),
		onlineAgentsChan: conf.OnlineAgentsChannel,
		retryTimes:       uint(conf.RetryTimes),
	}
}

func (c *CentrifugoClient) IsOnline(ctx context.Context, agentID string) (bool, error) {
	chans, err := c.ActiveChannels(ctx)
	if err != nil {
		return false, err
	}

	for _, ch := range chans {
		if ch == fmt.Sprintf("%s_message", agentID) {
			return true, nil
		}
	}
	return false, nil
}

func (c *CentrifugoClient) OnlineUsers(ctx context.Context, entID string) ([]string, error) {
	onlineUsers, err := c.ChannelUsers(ctx, fmt.Sprintf(c.onlineAgentsChan, entID))
	if err != nil {
		return nil, err
	}

	return onlineUsers, nil
}

func (c *CentrifugoClient) ChannelUsers(ctx context.Context, channel string) ([]string, error) {
	presenceResult, err := c.cli.Presence(ctx, channel)
	if err != nil {
		return nil, err
	}

	var users []string
	for _, client := range presenceResult.Presence {
		if client.User != "" {
			users = append(users, client.User)
		}
	}

	return users, nil
}

func (c *CentrifugoClient) PublishMessage(ctx context.Context, channel, message string) error {
	retryFunc := func() error {
		return c.cli.Publish(ctx, channel, []byte(message))
	}

	err := retry.Do(
		retryFunc,
		retry.Attempts(c.retryTimes),
		retry.Delay(5*time.Millisecond),
	)

	return err
}

func (c *CentrifugoClient) ActiveChannels(ctx context.Context) (chans []string, err error) {
	result, err := c.cli.Channels(ctx)
	if err != nil {
		return nil, err
	}
	return result.Channels, nil
}

func (c *CentrifugoClient) Unsubscribe(ctx context.Context, channel, user string) error {
	return c.cli.Unsubscribe(ctx, channel, user)
}
