package adapter

import (
	"time"

	"bitbucket.org/forfd/custm-chat/webim/common"
	"bitbucket.org/forfd/custm-chat/webim/common/logging"
	"bitbucket.org/forfd/custm-chat/webim/models"
)

// created_on: "2017-02-10T05:59:56.538108"
//enterprise_id: 5869
//id: 79093
//last_updated: "2019-01-19T08:19:56.968241"
//owner_id: 5869
//owner_type: "enterprise"
//quick_replies: [{content: "如有问题，可以拨打我们的客服电话 400 0031 322", content_type: "text",…},…]
//rank: 39050000
//title: "我们的联系方式
type QkReply struct {
	ID                 string    `json:"id"`
	EnterpriseID       string    `json:"enterprise_id"`
	GroupID            string    `json:"group_id"`
	Content            string    `json:"content"`
	ContentType        string    `json:"content_type"`
	CreatedOn          time.Time `json:"created_on"`
	HotKey             []int     `json:"hot_key"`
	KnowledgeConverted bool      `json:"knowledge_converted"`
	LastUpdated        time.Time `json:"last_updated"`
	Rank               int       `json:"rank"`
	RichContent        string    `json:"rich_content"`
	Title              string    `json:"title"`
}

type QkReplyGroup struct {
	ID           string     `json:"id"`
	EnterpriseID string     `json:"enterprise_id"`
	OwnerID      string     `json:"owner_id"`
	OwnerType    string     `json:"owner_type"`
	Rank         int        `json:"rank"`
	Title        string     `json:"title"`
	CreatedOn    time.Time  `json:"created_on"`
	LastUpdated  time.Time  `json:"last_updated"`
	QuickReplies []*QkReply `json:"quick_replies"`
}

type QkReplyGroupResp struct {
	AgentQuickReplyGroup []*QkReplyGroup `json:"agent_quick_reply_group"`
	EntQuickReplyGroup   []*QkReplyGroup `json:"ent_quick_reply_group"`
}

func ConvertModelQKReplyGroupToGroup(group *models.QuickreplyGroup) *QkReplyGroup {
	return &QkReplyGroup{
		ID:           group.ID,
		EnterpriseID: group.EntID,
		OwnerID:      group.CreatedBy,
		OwnerType:    group.CreatorType,
		Rank:         group.Rank,
		Title:        group.Title,
		CreatedOn:    group.CreatedAt,
		LastUpdated:  group.UpdatedAt,
	}
}

func ConvertModelQkItemToReply(entID string, item *models.QuickreplyItem) *QkReply {
	var hotKeys []int
	if err := common.Unmarshal(item.HotKey, &hotKeys); err != nil {
		log.Logger.Warnf("[ConvertModelQkItemToReply] unmarshal hotkeys error: %v", err)
	}

	return &QkReply{
		ID:                 item.ID,
		EnterpriseID:       entID,
		GroupID:            item.QuickreplyGroupID,
		Content:            item.Content,
		ContentType:        item.ContentType,
		CreatedOn:          item.CreatedAt,
		HotKey:             hotKeys,
		KnowledgeConverted: false,
		LastUpdated:        item.UpdatedAt,
		Rank:               item.Rank,
		RichContent:        item.RichContent.String,
		Title:              item.Title,
	}
}
