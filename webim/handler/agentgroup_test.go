package handler

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"bitbucket.org/forfd/custm-chat/webim/common"
	"bitbucket.org/forfd/custm-chat/webim/middleware"
	"bitbucket.org/forfd/custm-chat/webim/models"
	"bitbucket.org/forfd/custm-chat/webim/test"
)

type createAgentGroupResp struct {
	Code int                `json:"code"`
	Body *models.AgentGroup `json:"body"`
}

func TestIMService_CreateAgentGroup(t *testing.T) {
	test.InitTest()
	defer test.Clear()
	ast := assert.New(t)

	handler := &IMService{imCli: &test.IMClient{}, loc: test.IpLocation, mailClient: &test.MailClient{}}

	group := &CreateAgentGroupReq{
		Name:        "test group",
		Description: "test group desc",
	}
	v, _ := common.Marshal(group)
	_, rec, echoContext := setHttpRecorder(http.MethodPost, v)
	echoContext.Set(middleware.AgentEntIDKey, "1234")

	if ast.NoError(handler.CreateAgentGroup(echoContext)) {
		ast.Equal(http.StatusOK, rec.Code)

		var resp *createAgentGroupResp
		ast.NoError(common.Unmarshal(rec.Body.String(), &resp))
		ast.Equal(0, resp.Code)
		ast.Equal(group.Name, resp.Body.Name)
		ast.Equal(group.Description, resp.Body.Description)
		ast.Equal("1234", resp.Body.EntID)
	}

	group = &CreateAgentGroupReq{
		Name:        "",
		Description: "test group desc",
	}
	v, _ = common.Marshal(group)
	_, rec, echoContext = setHttpRecorder(http.MethodPost, v)
	echoContext.Set(middleware.AgentEntIDKey, "1234")

	if ast.NoError(handler.CreateAgentGroup(echoContext)) {
		ast.Equal(http.StatusOK, rec.Code)
		ast.Equal(`{"code":1000,"message":"name is invalid"}`, rec.Body.String())
	}

	group = &CreateAgentGroupReq{
		Name:        "test group",
		Description: "test group desc",
	}
	v, _ = common.Marshal(group)
	_, rec, echoContext = setHttpRecorder(http.MethodPost, v)
	echoContext.Set(middleware.AgentEntIDKey, "1234")

	if ast.NoError(handler.CreateAgentGroup(echoContext)) {
		ast.Equal(http.StatusOK, rec.Code)
		ast.Equal(`{"code":1000,"message":"name is exists"}`, rec.Body.String())
	}
}
