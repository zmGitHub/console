package elasticsearch

import (
	"context"

	"github.com/olivere/elastic"

	"bitbucket.org/forfd/custm-chat/webim/common/logging"
	"bitbucket.org/forfd/custm-chat/webim/conf"
)

var ESClient *esClient

type esClient struct {
	client *elastic.Client
}

type esLogger struct {
	level string
}

func (l *esLogger) Printf(format string, v ...interface{}) {
	switch l.level {
	case "info":
		log.Logger.Info(format, v)
	case "debug":
		log.Logger.Debug(format, v)
	case "error":
		log.Logger.Warnf(format, v)
	default:
	}
}

func NewESClient(config *conf.ElasticSearchConfig) error {
	options := []elastic.ClientOptionFunc{
		elastic.SetURL(config.Endpoint),
		elastic.SetSniff(false),
		elastic.SetErrorLog(&esLogger{level: "error"}),
		elastic.SetInfoLog(&esLogger{level: "info"}),
		elastic.SetTraceLog(&esLogger{level: "debug"}),
	}

	if len(config.Username) > 0 && len(config.Password) > 0 {
		options = append(options, elastic.SetBasicAuth(config.Username, config.Password))
	}

	client, err := elastic.NewClient(options...)
	if err != nil {
		return err
	}

	ESClient = &esClient{client: client}
	return nil
}

func (c *esClient) CreateDoc(index, itemType string, id string, doc interface{}) error {
	_, err := c.client.Index().
		Index(index).
		Type(itemType).
		Id(id).
		BodyJson(doc).
		Timeout(conf.IMConf.ElasticSearchConf.Timeout.String()).
		Refresh("wait_for").
		Do(context.Background())
	if err != nil {
		log.Logger.Warnf("ES CreateDoc error: %v", err)
		return err
	}

	return nil
}

func (c *esClient) CreateOrUpdateDoc(index string, itemType string, id string, doc map[string]interface{}) error {
	update := c.client.Update().
		Index(index).Type(itemType).Id(id).
		Doc(doc).
		DocAsUpsert(true).
		Timeout("1s").
		Refresh("true").
		DetectNoop(true)

	if _, err := update.Do(context.Background()); err != nil {
		log.Logger.Warnf("ES CreateOrUpdateDoc error: %v", err)
		return err
	}

	return nil
}

func (c *esClient) UpdateConversationEval(id string, level *int, content string) {
	c.UpdateConversation(id, map[string]interface{}{"eva_level": level, "eva_content": content})
}

func (c *esClient) UpdateConversationTags(id string, tags []string) {
	c.UpdateConversation(id, map[string]interface{}{"tags": tags})
}

func (c *esClient) UpdateConversationVisitorName(id string, visitorName string) {
	c.UpdateConversation(id, map[string]interface{}{"visitor.name": visitorName, "visit_info.name": visitorName})
}

func (c *esClient) UpdateConversation(id string, value map[string]interface{}) {
	update := c.client.Update().
		Index("conversation").Type("_doc").Id(id).
		Doc(value).
		Timeout("500ms").
		Refresh("wait_for").
		DetectNoop(true)

	if _, err := update.Do(context.Background()); err != nil {
		log.Logger.Warnf("ES UpdateConversation conversation_id=%s, error=%v", id, err)
	}
}
