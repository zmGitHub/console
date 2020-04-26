package handler

import (
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"bitbucket.org/forfd/custm-chat/webim/common"
	"bitbucket.org/forfd/custm-chat/webim/db"
	"bitbucket.org/forfd/custm-chat/webim/middleware"
	"bitbucket.org/forfd/custm-chat/webim/models"
	"bitbucket.org/forfd/custm-chat/webim/test"
)

type roleResp struct {
	Code int          `json:"code"`
	Body *models.Role `json:"body"`
}

type agentPermsResp struct {
	Code int            `json:"code"`
	Body []*models.Perm `json:"body"`
}

type rolePermsResp struct {
	Code int      `json:"code"`
	Body []string `json:"body"`
}

func TestIMService_AddRole(t *testing.T) {
	test.InitTest()
	defer test.Clear()
	ast := assert.New(t)

	handler := &IMService{imCli: &test.IMClient{}, loc: test.IpLocation, mailClient: &test.MailClient{}}

	data := &models.Role{
		Name: "admin",
	}

	v, _ := common.Marshal(data)
	_, rec, echoContext := setHttpRecorder(http.MethodPost, v)
	echoContext.Set(middleware.AgentEntIDKey, "1234")

	if ast.NoError(handler.AddRole(echoContext)) {
		var resp *roleResp
		ast.Equal(http.StatusOK, rec.Code)
		ast.NoError(common.Unmarshal(rec.Body.String(), &resp))
		ast.Equal(0, resp.Code)
		ast.Equal("1234", resp.Body.EntID)
		ast.Equal(data.Name, resp.Body.Name)
	}
}

func TestAddOrUpdatePermsToRole(t *testing.T) {
	test.InitTest()
	defer test.Clear()
	ast := assert.New(t)

	handler := &IMService{imCli: &test.IMClient{}, loc: test.IpLocation, mailClient: &test.MailClient{}}

	data := &AddPermsToRoleReq{
		PermIDs: []string{"1a2b", "3a4d"},
	}

	v, _ := common.Marshal(data)
	_, rec, echoContext := setHttpRecorder(http.MethodPost, v)
	echoContext.SetPath("/admin/api/v1/roles/:role_id/perms")
	echoContext.SetParamNames("role_id")
	echoContext.SetParamValues("5678")

	if ast.NoError(handler.AddOrUpdatePermsToRole(echoContext)) {
		ast.Equal(http.StatusOK, rec.Code)
		ast.Equal(`{"code":0,"body":{"perm_ids":["1a2b","3a4d"]}}`, rec.Body.String())

		perms, err := models.RolePermsByRoleID(db.Mysql, "5678")
		ast.Nil(err)
		ast.Equal(2, len(perms))
		for _, p := range perms {
			log.Println(p.RoleID, p.PermID)
		}
	}
}

func TestIMService_GetAgentPerms(t *testing.T) {
	test.InitTest()
	defer test.Clear()
	ast := assert.New(t)

	handler := &IMService{imCli: &test.IMClient{}, loc: test.IpLocation, mailClient: &test.MailClient{}}

	now := time.Now().UTC()
	agent := &models.Agent{
		ID:            common.GenUniqueID(),
		EntID:         "1234",
		GroupID:       "",
		RoleID:        "5678",
		AccountStatus: models.AgentAccountCreatedStatus,
		CreateAt:      now,
		UpdateAt:      now,
	}
	ast.NoError(agent.Insert(db.Mysql))
	perms := []*models.Perm{
		{
			ID:        common.GenUniqueID(),
			EntID:     "1234",
			AppName:   "conversation",
			Name:      "check",
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			ID:        common.GenUniqueID(),
			EntID:     "1234",
			AppName:   "report",
			Name:      "check",
			CreatedAt: now,
			UpdatedAt: now,
		},
	}
	for _, p := range perms {
		ast.NoError(p.Insert(db.Mysql))
	}

	roles := []*models.RolePerm{
		{
			RoleID:    "5678",
			PermID:    perms[0].ID,
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			RoleID:    "5678",
			PermID:    perms[1].ID,
			CreatedAt: now,
			UpdatedAt: now,
		},
	}
	for _, r := range roles {
		ast.NoError(r.Insert(db.Mysql))
	}

	rec, echoContext := setHttpGetCtx()
	echoContext.SetPath("/admin/api/v1/agents/:agent_id/perms")
	echoContext.SetParamNames("agent_id")
	echoContext.SetParamValues(agent.ID)

	if ast.NoError(handler.GetAgentPerms(echoContext)) {
		var resp *agentPermsResp
		ast.Equal(http.StatusOK, rec.Code)
		ast.NoError(common.Unmarshal(rec.Body.String(), &resp))
		ast.Equal(0, resp.Code)
		ast.Equal(2, len(resp.Body))
		for i, perm := range resp.Body {
			ast.Equal("1234", perm.EntID)
			ast.Equal(perms[i].AppName, perm.AppName)
			ast.Equal(perms[i].Name, perm.Name)
		}
	}
}

func TestIMService_GetRolePerms(t *testing.T) {
	test.InitTest()
	defer test.Clear()
	ast := assert.New(t)

	handler := &IMService{imCli: &test.IMClient{}, loc: test.IpLocation, mailClient: &test.MailClient{}}

	now := time.Now().UTC()
	perms := []*models.Perm{
		{
			ID:        common.GenUniqueID(),
			EntID:     "1234",
			AppName:   "conversation",
			Name:      "check",
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			ID:        common.GenUniqueID(),
			EntID:     "1234",
			AppName:   "report",
			Name:      "check",
			CreatedAt: now,
			UpdatedAt: now,
		},
	}
	for _, p := range perms {
		ast.NoError(p.Insert(db.Mysql))
	}

	roles := []*models.RolePerm{
		{
			RoleID:    "5678",
			PermID:    perms[0].ID,
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			RoleID:    "5678",
			PermID:    perms[1].ID,
			CreatedAt: now,
			UpdatedAt: now,
		},
	}
	for _, r := range roles {
		ast.NoError(r.Insert(db.Mysql))
	}

	rec, echoContext := setHttpGetCtx()
	echoContext.SetPath("/admin/api/v1/roles/:role_id/perms")
	echoContext.SetParamNames("role_id")
	echoContext.SetParamValues("5678")

	if ast.NoError(handler.GetRolePerms(echoContext)) {
		var resp *rolePermsResp
		ast.Equal(http.StatusOK, rec.Code)
		ast.NoError(common.Unmarshal(rec.Body.String(), &resp))
		ast.Equal(0, resp.Code)
		ast.Equal(2, len(resp.Body))
		ast.Equal([]string{perms[0].ID, perms[1].ID}, resp.Body)
	}
}
