package imclient

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"bitbucket.org/forfd/custm-chat/webim/conf"
)

func TestCentrifugoClient_PublishMessage(t *testing.T) {
	client := NewClient(&conf.CentrifugoConfig{
		API:        "http://localhost:9000/api",
		AuthKey:    "",
		Timeout:    conf.Duration{Duration: 2 * time.Second},
		RetryTimes: 3,
	})

	err := client.PublishMessage(context.Background(), "news", "news come")
	assert.Nil(t, err)
}
