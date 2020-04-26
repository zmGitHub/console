package handler

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"bitbucket.org/forfd/custm-chat/webim/common"
	"bitbucket.org/forfd/custm-chat/webim/db"
	"bitbucket.org/forfd/custm-chat/webim/middleware"
	"bitbucket.org/forfd/custm-chat/webim/models"
	"bitbucket.org/forfd/custm-chat/webim/test"
)

func setHttpGetCtx() (*httptest.ResponseRecorder, echo.Context) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	echoContext := e.NewContext(req, rec)
	return rec, echoContext
}

func TestIMService_GetEntAgents(t *testing.T) {
	test.InitTest()
	defer test.Clear()

	ast := assert.New(t)
	handler := &IMService{imCli: &test.IMClient{}, loc: test.IpLocation, mailClient: &test.MailClient{}}
	rec, echoContext := setHttpGetCtx()
	echoContext.Set("EntID", "xxxxx")

	if ast.NoError(handler.GetEntAgents(echoContext)) {
		ast.Equal(http.StatusOK, rec.Code)
		ast.Equal(`{"code":0,"body":null}`, rec.Body.String())
	}
}

func TestIMService_GetAgentByID(t *testing.T) {
	test.InitTest()
	defer test.Clear()

	ast := assert.New(t)
	handler := &IMService{imCli: &test.IMClient{}, loc: test.IpLocation, mailClient: &test.MailClient{}}
	rec, echoContext := setHttpGetCtx()
	echoContext.SetPath("/admin/api/v1/agents/:agent_id")
	echoContext.SetParamNames("agent_id")
	echoContext.SetParamValues("xxxxxx")

	if ast.NoError(handler.GetAgentByID(echoContext)) {
		ast.Equal(http.StatusOK, rec.Code)
		ast.Equal(`{"code":1003,"message":"agent not exists"}`, rec.Body.String())
	}

	rec, echoContext = setHttpGetCtx()
	echoContext.SetPath("/admin/api/v1/agents/:agent_id")
	echoContext.SetParamNames("agent_id")
	echoContext.SetParamValues("")

	if ast.NoError(handler.GetAgentByID(echoContext)) {
		ast.Equal(http.StatusOK, rec.Code)
		ast.Equal(`{"code":1000,"message":"agent_id is invalid"}`, rec.Body.String())
	}
}

func TestIMService_GetEntAgentGroups(t *testing.T) {
	groups := []*models.AgentGroup{
		{ID: common.GenUniqueID(), EntID: "1234", Name: "group1", Description: "gp1"},
		{ID: common.GenUniqueID(), EntID: "1234", Name: "group2", Description: "gp2"},
		{ID: common.GenUniqueID(), EntID: "1234", Name: "group3", Description: "gp3"},
	}

	test.InitTest()
	defer test.Clear()
	ast := assert.New(t)

	handler := &IMService{imCli: &test.IMClient{}, loc: test.IpLocation, mailClient: &test.MailClient{}}
	rec, echoContext := setHttpGetCtx()
	echoContext.SetPath("/admin/api/v1/enterprises/agent_groups")
	echoContext.Set(middleware.AgentEntIDKey, "1234")

	for _, group := range groups {
		ast.Nil(group.Insert(db.Mysql))
	}

	type getGroupsResp struct {
		Code   int                  `json:"code"`
		Groups []*models.AgentGroup `json:"body"`
	}
	if ast.NoError(handler.GetEntAgentGroups(echoContext)) {
		ast.Equal(http.StatusOK, rec.Code)

		var resp *getGroupsResp
		err := common.Unmarshal(rec.Body.String(), &resp)
		ast.Nil(err)
		ast.Equal(0, resp.Code)
		for i, group := range resp.Groups {
			ast.Equal(groups[i].ID, group.ID)
			ast.Equal(groups[i].Name, group.Name)
			ast.Equal(groups[i].EntID, group.EntID)
			ast.Equal(groups[i].Description, group.Description)
			log.Println(group.ID, group.Name, group.EntID, group.Description)
		}
	}
}

func TestIMService_UpdateAgentStatus(t *testing.T) {
	test.InitTest()
	defer test.Clear()
	ast := assert.New(t)

	handler := &IMService{imCli: &test.IMClient{}, loc: test.IpLocation, mailClient: &test.MailClient{}}
	cases := []struct {
		data       string
		expectResp string
	}{
		{
			data:       `{"status":""}`,
			expectResp: `{"code":1000,"message":"unsupported status"}`,
		},
		{
			data:       `{"status":"online"}`,
			expectResp: `{"code":0,"body":null}`,
		},
		{
			data:       `{"status":"offline"}`,
			expectResp: `{"code":0,"body":null}`,
		},
	}
	for _, c := range cases {
		_, rec, echoContext := setHttpRecorder(http.MethodPut, c.data)
		echoContext.Set(middleware.AgentEntIDKey, "1234")
		echoContext.Set(middleware.AgentIDKey, "5678")

		if ast.NoError(handler.UpdateAgentStatus(echoContext)) {
			ast.Equal(http.StatusOK, rec.Code)
			ast.Equal(c.expectResp, rec.Body.String())
		}
	}
}
