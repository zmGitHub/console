package dto

type SelectingRuleTarget []interface{}

type SelectingRule struct {
	CreatedOn    string                `json:"created_on"`
	EnterpriseID string                `json:"enterprise_id"`
	ID           string                `json:"id"`
	Inverted     bool                  `json:"inverted"`
	LastUpdated  string                `json:"last_updated"`
	MatchRules   [][]string            `json:"match_rules"`
	Rank         int                   `json:"rank"`
	Targets      []SelectingRuleTarget `json:"targets"`
	Type         string                `json:"type"`
}
