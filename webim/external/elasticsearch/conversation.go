package elasticsearch

import (
	"context"
	"reflect"

	"github.com/olivere/elastic"

	"bitbucket.org/forfd/custm-chat/webim/common/logging"
	"bitbucket.org/forfd/custm-chat/webim/handler/adapter"
)

type Message struct {
	ID      string `json:"id"`
	Content string `json:"content"`
}

type SimpleVisitor struct {
	Name      string `json:"name"`
	Telephone string `json:"telephone"`
	Remark    string `json:"remark"`
}

type ConversationV1 struct {
	*adapter.Conversation
	Visitor  *SimpleVisitor `json:"visitor"`
	Messages []*Message     `json:"messages"`
}

// CreateConversationDoc ...
func CreateConversationDoc(esCli *esClient, doc *ConversationV1) error {
	return esCli.CreateDoc("conversation", "_doc", doc.ID, doc)
}

// SearchConversations search conversations
func SearchConversations(esCli *esClient, query elastic.Query, from, size int) (res []*ConversationV1, totalCount int64, err error) {
	searchResult, err := esCli.client.Search().
		Index("conversation").
		Type("_doc").
		Query(query).
		Sort("created_on", false).
		From(from).Size(size).
		Pretty(true).
		Do(context.Background())
	if err != nil {
		var rootReason string
		if searchResult != nil && len(searchResult.Error.RootCause) > 0 {
			rootReason = searchResult.Error.RootCause[0].Reason
		}

		log.Logger.Warnf("[SearchConversations] root reason: %s, error: %v\n", rootReason, err)
		return nil, 0, err
	}

	if searchResult == nil {
		return []*ConversationV1{}, 0, nil
	}

	for _, hit := range searchResult.Each(reflect.TypeOf(&ConversationV1{})) {
		v, ok := hit.(*ConversationV1)
		if ok {
			res = append(res, v)
		}
	}

	return res, searchResult.Hits.TotalHits, err
}

func UpdateConversationEvalLevel(esCli *esClient, conversationID string, level *int, content string) {
	esCli.UpdateConversationEval(conversationID, level, content)
}

func UpdateConversationTags(esCli *esClient, conversationID string, tags []string) {
	esCli.UpdateConversationTags(conversationID, tags)
}

func UpdateConversationVisitorName(esCli *esClient, conversationID string, visitorName string) {
	esCli.UpdateConversationVisitorName(conversationID, visitorName)
}
