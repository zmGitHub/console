package adapter

type Attribute struct {
	Deleted      bool `json:"deleted"`
	EnterpriseID int  `json:"enterprise_id"`
	ID           int  `json:"id"`
	MetaData     []struct {
		Deleted bool   `json:"deleted"`
		Name    string `json:"name"`
		Value   string `json:"value"`
	} `json:"meta_data"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type CustomFields struct {
	CustomAttrs []*Attribute `json:"custom_attrs"`
}

type AttrOrder struct {
	Orders []struct {
		AttrName string `json:"attr_name"`
		Visible  bool   `json:"visible"`
	} `json:"orders"`
}
