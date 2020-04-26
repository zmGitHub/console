package test

import (
	"context"
	"log"
)

type IMClient struct {
	OnlineUsers []string
}

func (c *IMClient) IsOnline(ctx context.Context, entID string) (bool, error) {
	return true, nil
}

func (c *IMClient) ActiveChannels(ctx context.Context) (chans []string, err error) {
	return nil, nil
}

func (c *IMClient) ChannelUsers(ctx context.Context, channel string) ([]string, error) {
	return c.OnlineUsers, nil
}

func (c *IMClient) PublishMessage(ctx context.Context, channel, message string) error {
	log.Println(channel, message)
	return nil
}

func (c *IMClient) Unsubscribe(ctx context.Context, channel, user string) error {
	return nil
}
