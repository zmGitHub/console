package handler

type ErrMsg struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type Resp struct {
	Code int         `json:"code"`
	Body interface{} `json:"body,omitempty"`
}

type SuccessResp struct {
	Success bool `json:"success"`
}
