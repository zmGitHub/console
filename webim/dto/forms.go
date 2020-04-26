package dto

// {
//          "description": "其他",
//          "target": "",
//          "target_kind": ""
//        }
var DefaultChatForms = &ChatForms{
	FormDef: &FormDef{
		Captcha: false,
		Inputs: &Inputs{
			Fields: []*Field{},
			Status: "close",
		},
		Menus: &Menus{
			Assignments: []*Assignment{
				{
					Description: "其他",
					TargetKind:  "",
					Target:      "",
				},
			},
			Title:  "问题",
			Status: "close",
		},
	},
}

type Assignment struct {
	Description string `json:"description"`
	Target      string `json:"target"`
	TargetKind  string `json:"target_kind"`
}

type Menus struct {
	Assignments []*Assignment `json:"assignments"`
	Status      string        `json:"status"`
	Title       string        `json:"title"`
}

type Field struct {
	DisplayName            string   `json:"display_name"`
	FieldName              string   `json:"field_name"`
	IgnoreReturnedCustomer bool     `json:"ignore_returned_customer"`
	Optional               bool     `json:"optional"`
	Type                   string   `json:"type"`
	Choices                []string `json:"choices"`
}

type Inputs struct {
	Fields []*Field `json:"fields"`
	Status string   `json:"status"`
	Title  string   `json:"title"`
}

type FormDef struct {
	Captcha bool    `json:"captcha"`
	Inputs  *Inputs `json:"inputs"`
	Menus   *Menus  `json:"menus"`
	Version int     `json:"version"`
}

type ChatForms struct {
	ID           string   `json:"id"`
	EnterpriseID string   `json:"enterprise_id"`
	Title        string   `json:"title"`
	FormDef      *FormDef `json:"form_def"`
	CreatedOn    string   `json:"created_on"`
	LastUpdated  string   `json:"last_updated"`
}
