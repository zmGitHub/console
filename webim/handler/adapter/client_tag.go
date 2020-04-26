package adapter

import (
	"time"

	"bitbucket.org/forfd/custm-chat/webim/models"
)

// "color": "tag-seagreen",
//    "created_on": "2019-03-31T00:05:52.250583",
//    "creator_id": 9422,
//    "enterprise_id": 5869,
//    "id": 33431,
//    "last_updated": "2019-03-31T00:05:52.250591",
//    "name": "new_tag111",
//    "rank": 810000,
//    "use_count": 0
type ClientTag struct {
	ID           string    `json:"id"`
	Color        string    `json:"color"`
	CreatedOn    time.Time `json:"created_on"`
	CreatorID    string    `json:"creator_id"`
	EnterpriseID string    `json:"enterprise_id"`
	LastUpdated  time.Time `json:"last_updated"`
	Name         string    `json:"name"`
	Rank         int       `json:"rank"`
	UseCount     int       `json:"use_count"`
}

type ClientTagsResp struct {
	ClientTags []*ClientTag `json:"client_tags"`
}

func ConvertModelTagToClientTag(tag *models.VisitorTag) *ClientTag {
	return &ClientTag{
		ID:           tag.ID,
		Color:        tag.Color,
		CreatedOn:    tag.CreatedAt,
		CreatorID:    tag.Creator,
		EnterpriseID: tag.EntID,
		LastUpdated:  tag.UpdatedAt,
		Name:         tag.Name,
		Rank:         tag.Rank,
		UseCount:     tag.UseCount,
	}
}

func ConvertVisitorTagsToClientTags(tags []*models.VisitorTag) *ClientTagsResp {
	resp := &ClientTagsResp{
		ClientTags: make([]*ClientTag, len(tags)),
	}

	for i, tag := range tags {
		resp.ClientTags[i] = ConvertModelTagToClientTag(tag)
	}

	return resp
}
