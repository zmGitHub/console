package handler

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"bitbucket.org/forfd/custm-chat/webim/common"
	"bitbucket.org/forfd/custm-chat/webim/middleware"
	"bitbucket.org/forfd/custm-chat/webim/test"
)

func TestIMService_CreateEvaluation(t *testing.T) {
	test.InitTest()
	defer test.Clear()
	ast := assert.New(t)

	handler := &IMService{imCli: &test.IMClient{}, loc: test.IpLocation, mailClient: &test.MailClient{}}

	data := &evaluation{
		AgentID: "1234",
		Level:   5,
		Content: "test eval content",
	}
	v, _ := common.Marshal(data)
	_, rec, echoContext := setHttpRecorder(http.MethodPost, v)
	echoContext.Set(middleware.AgentEntIDKey, "5678")
	echoContext.SetPath("/api/v1/enterprises/:ent_id/evaluations")
	echoContext.SetParamNames("ent_id")
	echoContext.SetParamValues("5678")

	type resp struct {
		Code int             `json:"code"`
		Body *evaluationResp `json:"body"`
	}
	if ast.NoError(handler.CreateEvaluation(echoContext)) {
		ast.Equal(http.StatusOK, rec.Code)

		var res *resp
		ast.NoError(common.Unmarshal(rec.Body.String(), &res))
		ast.Equal(0, res.Code)
		ast.Equal(5, res.Body.Level)
	}
}
