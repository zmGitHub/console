package handler

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/labstack/echo/v4"

	"bitbucket.org/forfd/custm-chat/webim/common"
	log "bitbucket.org/forfd/custm-chat/webim/common/logging"
	"bitbucket.org/forfd/custm-chat/webim/conf"
	"bitbucket.org/forfd/custm-chat/webim/db"
	"bitbucket.org/forfd/custm-chat/webim/dto"
	"bitbucket.org/forfd/custm-chat/webim/middleware"
	"bitbucket.org/forfd/custm-chat/webim/models"
)

type formField struct {
	DisplayName string `json:"display_name"`
	FieldName   string `json:"field_name"`
	ValueType   string `json:"value_type"`
	Required    bool   `json:"required"`
}

func (f *formField) validate() error {
	if f.DisplayName == "" {
		return fmt.Errorf("display_name is empty")
	}

	if f.FieldName == "" {
		return fmt.Errorf("field_name is empty")
	}

	if f.ValueType == "" {
		return fmt.Errorf("value_type is empty")
	}

	return nil
}

type menuField struct {
	Description string `json:"description"`
	Value       string `json:"value"`
	AgentType   string `json:"agent_type"`
}

func (m *menuField) validate() error {
	if m.Description == "" {
		return fmt.Errorf("menu description is empty")
	}

	if m.Value == "" {
		return fmt.Errorf("menu value is empty")
	}

	if m.AgentType != "agent" && m.AgentType != "agent_group" {
		return fmt.Errorf("unsupported menu agent type")
	}

	return nil
}

type inputs struct {
	Status      string       `json:"status"`
	Description string       `json:"description"`
	Fields      []*formField `json:"fields"`
}

func (i *inputs) validate() error {
	if i.Status == "" {
		return fmt.Errorf("inputs status is empty")
	}

	if i.Status != "open" && i.Status != "close" {
		return fmt.Errorf("unsupported inputs status")
	}

	if i.Description == "" {
		return fmt.Errorf("inputs description is empty")
	}

	if i.Fields == nil {
		return fmt.Errorf("fields is nil")
	}

	for _, field := range i.Fields {
		err := field.validate()
		if err != nil {
			return err
		}
	}

	return nil
}

type menu struct {
	Status      string       `json:"status"`
	Description string       `json:"description"`
	Fields      []*menuField `json:"fields"`
}

func (m *menu) validate() error {
	if m.Status == "" {
		return fmt.Errorf("inputs status is empty")
	}

	if m.Status != "open" && m.Status != "close" {
		return fmt.Errorf("unsupported inputs status ")
	}

	if m.Description == "" {
		return fmt.Errorf("inputs description is empty")
	}

	if m.Fields == nil {
		return fmt.Errorf("fields is nil")
	}

	for _, field := range m.Fields {
		err := field.validate()
		if err != nil {
			return err
		}
	}

	return nil
}

type PrechatFormFields struct {
	Inputs *inputs `json:"inputs"`
	Menus  *menu   `json:"menus"`
}

type prechatForm struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	*PrechatFormFields
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (form *prechatForm) validate() error {
	if form.Title == "" {
		return fmt.Errorf("title is empty")
	}

	if form.Inputs == nil {
		return fmt.Errorf("form inputs is nil")
	}

	if err := form.Inputs.validate(); err != nil {
		return err
	}

	if form.Menus == nil {
		return fmt.Errorf("form menus is nil")
	}

	if err := form.Menus.validate(); err != nil {
		return err
	}

	return nil
}

type FillFormRequest struct {
	EntID   string  `json:"ent_id"`
	TrackID string  `json:"track_id"`
	Data    attrsV1 `json:"data"`
}

// AddForm ...
// POST /admin/api/v1/prechat_forms
func (s *IMService) AddForm(ctx echo.Context) (err error) {
	entID := ctx.Get(middleware.AgentEntIDKey).(string)
	agentID := ctx.Get(middleware.AgentIDKey).(string)

	if !conf.IMConf.Debug {
		if msg := hasPerm(entID, agentID, "agent_settings", "check_update_prechat_form"); msg != nil {
			return noPermResp(ctx, msg)
		}
	}

	form := &dto.ChatForms{}
	if err = ctx.Bind(form); err != nil {
		return invalidParameterResp(ctx, err.Error())
	}

	formID, err := s.addForm(entID, form.Title, db.Mysql, form.FormDef)
	if err != nil {
		return invalidParameterResp(ctx, err.Error())
	}

	form.ID = formID
	return jsonResponse(ctx, form)
}

func (s *IMService) addForm(entID, title string, tx models.XODB, def *dto.FormDef) (id string, err error) {
	fields, err := common.Marshal(def)
	if err != nil {
		return "", err
	}

	now := time.Now().UTC()
	chatForm := &models.PrechatForm{
		ID:         common.GenUniqueID(),
		EntID:      entID,
		Title:      title,
		FormFields: sql.NullString{String: fields, Valid: true},
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	if err := chatForm.Insert(tx); err != nil {
		return "", err
	}

	return chatForm.ID, nil
}

// UpdateForm ...
// PUT /api/agent/forms/:form_id
func (s *IMService) UpdateForm(ctx echo.Context) (err error) {
	entID := ctx.Get(middleware.AgentEntIDKey).(string)
	agentID := ctx.Get(middleware.AgentIDKey).(string)

	if !conf.IMConf.Debug {
		if msg := hasPerm(entID, agentID, "online_agent_config", "config_prechat_survey"); msg != nil {
			return noPermResp(ctx, msg)
		}
	}

	formID := ctx.Param("form_id")
	if formID == "" {
		return invalidParameterResp(ctx, "form_id is empty")
	}

	form := &dto.ChatForms{}
	if err = ctx.Bind(form); err != nil {
		return
	}

	fields, err := common.Marshal(form.FormDef)
	if err != nil {
		return errResp(ctx, common.EncodeJSONErr, err.Error())
	}

	chatForm, err := models.PrechatFormByID(db.Mysql, formID)
	if err != nil {
		return dbErrResp(ctx, err.Error())
	}

	chatForm.FormFields = sql.NullString{String: fields, Valid: true}
	chatForm.UpdatedAt = time.Now().UTC()
	if err = chatForm.Update(db.Mysql); err != nil {
		return dbErrResp(ctx, err.Error())
	}

	if form.FormDef != nil && form.FormDef.Inputs != nil {
		config, errMsg := s.getEnterpriseConfigs(entID)
		if form.FormDef.Inputs.Status == "open" || form.FormDef.Menus.Status == "open" {
			if errMsg == nil && config != nil {
				if config.Survey != nil {
					config.Survey.Status = "open"
					s.updateEntConfigs(entID, config)
				}
			}
		} else {
			if errMsg == nil && config != nil {
				if config.Survey != nil {
					config.Survey.Status = "close"
					s.updateEntConfigs(entID, config)
				}
			}
		}
	}

	form.ID = chatForm.ID
	form.EnterpriseID = chatForm.EntID
	form.CreatedOn = *common.ConvertUTCToTimeString(chatForm.CreatedAt)
	form.LastUpdated = *common.ConvertUTCToTimeString(chatForm.UpdatedAt)

	return jsonResponse(ctx, form)
}

// GET /client/forms?ent_id=xxx
func (s *IMService) GetClientForms(ctx echo.Context) (err error) {
	entID := ctx.QueryParam("ent_id")
	if entID == "" {
		return invalidParameterResp(ctx, "Invalid ent_id")
	}

	res, err := s.getEntForms(entID)
	if err != nil {
		return invalidParameterResp(ctx, err.Error())
	}

	return jsonResponse(ctx, res)
}

// POST /client/forms
// /client/forms?ent_id=xxxx
func (s *IMService) FillForms(ctx echo.Context) (err error) {
	req := &FillFormRequest{}
	if err = ctx.Bind(req); err != nil {
		return invalidParameterResp(ctx, err.Error())
	}

	visitor, err := models.VisitorByEntIDTraceID(db.Mysql, req.EntID, req.TrackID)
	if err != nil {
		if err == sql.ErrNoRows {
			return invalidParameterResp(ctx, "visitor not found")
		}

		return dbErrResp(ctx, err.Error())
	}

	// NewAttrsChange
	var attrsValue = map[string]interface{}{}
	if err := s.updateVisitorAttrs(req.EntID, req.TrackID, visitor, req.Data, attrsValue); err != nil {
		log.Logger.Warnf("error: %v", err)
	}

	return jsonResponse(ctx, &SuccessResp{Success: true})
}

// GET /api/agent/forms
func (s *IMService) GetEntForms(ctx echo.Context) (err error) {
	entID := ctx.Get(middleware.AgentEntIDKey).(string)
	agentID := ctx.Get(middleware.AgentIDKey).(string)

	if !conf.IMConf.Debug {
		if msg := hasPerm(entID, agentID, "online_agent_config", "config_prechat_survey"); msg != nil {
			return noPermResp(ctx, msg)
		}
	}
	res, err := s.getEntForms(entID)
	if err != nil {
		return invalidParameterResp(ctx, err.Error())
	}

	return jsonResponse(ctx, res)
}

func (s *IMService) getEntForms(entID string) (forms *dto.ChatForms, err error) {
	modelForm, err := models.PrechatFormsByEntID(db.Mysql, entID)
	if err != nil {
		if err == sql.ErrNoRows {
			return dto.DefaultChatForms, nil
		}

		return nil, err
	}

	formStr := `{}`
	if modelForm.FormFields.Valid {
		formStr = modelForm.FormFields.String
	}

	res := &dto.ChatForms{
		ID:           modelForm.ID,
		EnterpriseID: modelForm.EntID,
		Title:        modelForm.Title,
		CreatedOn:    *common.ConvertUTCToTimeString(modelForm.CreatedAt),
		LastUpdated:  *common.ConvertUTCToTimeString(modelForm.UpdatedAt),
	}

	if err = common.Unmarshal(formStr, &res.FormDef); err != nil {
		return res, err
	}

	return res, nil
}
