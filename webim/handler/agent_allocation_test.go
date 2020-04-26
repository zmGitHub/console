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

func TestIMService_AllocateAgent(t *testing.T) {

}

func TestIMService_CreateOrUpdateAllocationRule(t *testing.T) {
	test.InitTest()
	defer test.Clear()
	ast := assert.New(t)

	handler := &IMService{imCli: &test.IMClient{}, loc: test.IpLocation, mailClient: &test.MailClient{}}
	cases := []struct {
		data       string
		expectResp string
	}{
		{
			data:       `{"rule_type":""}`,
			expectResp: `{"code":1000,"message":"unsupported rule_type"}`,
		},
		{
			data:       `{"rule_type":"load_balanced"}`,
			expectResp: `{"code":1000,"message":"unsupported rule_type"}`,
		},
	}
	for _, c := range cases {
		_, rec, echoContext := setHttpRecorder(http.MethodPost, c.data)
		echoContext.Set(middleware.AgentEntIDKey, "1234")

		if ast.NoError(handler.CreateOrUpdateAllocationRule(echoContext)) {
			ast.Equal(http.StatusOK, rec.Code)
			ast.Equal(c.expectResp, rec.Body.String())
		}
	}

	data := `{"rule_type":"conversation_num"}`
	_, rec, echoContext := setHttpRecorder(http.MethodPost, data)
	echoContext.Set(middleware.AgentEntIDKey, "1234")

	if ast.NoError(handler.CreateOrUpdateAllocationRule(echoContext)) {
		ast.Equal(http.StatusOK, rec.Code)
		var resp = struct {
			Code int                    `json:"code"`
			Body *models.AllocationRule `json:"body"`
		}{}
		err := common.Unmarshal(rec.Body.String(), &resp)
		ast.Nil(err)
		ast.Equal(0, resp.Code)
		ast.Equal("1234", resp.Body.EntID)
		ast.Equal("conversation_num", resp.Body.RuleType)
	}
}
