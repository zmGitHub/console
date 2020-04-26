package adapter

type TicketCategory struct {
	ID       int    `json:"id"`
	IsFolder int    `json:"is_folder"`
	Name     string `json:"name"`
	ParentID int    `json:"parent_id"`
	Subcat   []struct {
		ID       int           `json:"id"`
		IsFolder int           `json:"is_folder"`
		Name     string        `json:"name"`
		ParentID int           `json:"parent_id"`
		Subcat   []interface{} `json:"subcat"`
	} `json:"subcat"`
}

type Tickets struct {
	Categories []*TicketCategory `json:"categories"`
}
