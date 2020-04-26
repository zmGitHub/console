package dto

import (
	"bitbucket.org/forfd/custm-chat/webim/common"
	"bitbucket.org/forfd/custm-chat/webim/models"
)

// {"name":" 昌平","rules":[{"attribute":"country","condition":"eq","type":"string","value":"中国"},{"attribute":"address","condition":"contain","type":"string","value":"北京"}]}

type SegmentsRule struct {
	Attribute string `json:"attribute"`
	Condition string `json:"condition"`
	Type      string `json:"type"`
	Value     string `json:"value"`
}

type Segments struct {
	Name  string          `json:"name"`
	Rules []*SegmentsRule `json:"rules"`
}

type UpdateSegmentsReq struct {
	Name string `json:"name"`
}

type CreateSegmentsResp struct {
	ID           string `json:"id"`
	CreatedOn    string `json:"created_on"`
	UpdatedOn    string `json:"updated_on"`
	EnterpriseID string `json:"enterprise_id"`
	*Segments
}

type GetSegmentsResp struct {
	Segments []*CreateSegmentsResp `json:"segments"`
}

func ConvertToSegmentsResp(segment *models.UserSegment) (*CreateSegmentsResp, error) {
	s := &Segments{
		Name: segment.Name,
	}

	if err := common.Unmarshal(segment.Rules.String, &s.Rules); err != nil {
		return nil, err
	}

	return &CreateSegmentsResp{
		ID:           segment.ID,
		CreatedOn:    *common.ConvertUTCToTimeString(segment.CreatedAt),
		UpdatedOn:    *common.ConvertUTCToTimeString(segment.CreatedAt),
		EnterpriseID: segment.EntID,
		Segments:     s,
	}, nil
}
