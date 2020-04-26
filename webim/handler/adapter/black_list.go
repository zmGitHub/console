package adapter

type VisitBlackListResp struct {
	Success        bool          `json:"success"`
	TotalCount     int           `json:"total_count"`
	VisitBlacklist []interface{} `json:"visit_blacklist"`
}
