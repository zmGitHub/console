package imclient

import (
	"context"
)

type IMClient interface {
	IsOnline(ctx context.Context, agentID string) (bool, error)
	ChannelUsers(ctx context.Context, channel string) ([]string, error)
	PublishMessage(ctx context.Context, channel, message string) error
	ActiveChannels(ctx context.Context) (chans []string, err error)
	Unsubscribe(ctx context.Context, channel, user string) error
}
