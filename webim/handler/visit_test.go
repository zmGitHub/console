package handler

import (
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"bitbucket.org/forfd/custm-chat/webim/common"
	"bitbucket.org/forfd/custm-chat/webim/db"
	"bitbucket.org/forfd/custm-chat/webim/models"
	"bitbucket.org/forfd/custm-chat/webim/test"
)

type initVisitResp struct {
	Code int            `json:"code"`
	Body *InitVisitResp `json:"body"`
}

func setHttpRecorder(method string, data string) (req *http.Request, rec *httptest.ResponseRecorder, echoContext echo.Context) {
	e := echo.New()
	req = httptest.NewRequest(method, "/", strings.NewReader(data))
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.102 Safari/537.36")
	req.Header.Set(echo.HeaderXRealIP, "123.125.115.110") // // X-Real-Ip
	req.Header.Set("Referer", "https://test.custmchat.io")
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	echoContext = e.NewContext(req, rec)
	return
}

func TestIMService_InitVisitNoTrace(t *testing.T) {
	test.InitTest()
	defer test.Clear()

	ast := assert.New(t)
	handler := &IMService{loc: test.IpLocation}
	data := &InitVisitReq{
		EntID:   common.GenUniqueID(),
		TraceID: "",
		Keyword: "test key word",
		Title:   "test title",
	}

	v, err := common.Marshal(data)
	ast.Nil(err)

	_, rec, echoContext := setHttpRecorder(http.MethodPost, v)
	echoContext.SetPath("/api/v1/enterprises/:ent_id/visits")
	echoContext.SetParamNames("ent_id")
	echoContext.SetParamValues(data.EntID)

	if ast.NoError(handler.InitVisit(echoContext)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		body := rec.Body.String()

		var s *initVisitResp
		err = common.Unmarshal(body, &s)
		ast.Nil(err)

		ast.Equal(s.Code, 0)
		modelVisit, err := models.VisitByID(db.Mysql, s.Body.VisitID)
		modelVisits, err := models.VisitsByEntIDTraceID(db.Mysql, data.EntID, s.Body.TraceID)
		ast.Nil(err)
		ast.Equal(data.EntID, modelVisit.EntID)
		ast.Equal(modelVisit, modelVisits[0])
		log.Println(
			"os family: ", modelVisit.OsFamily,
			"category: ", modelVisit.OsCategory,
			"os version: ", modelVisit.OsVersion,
			"os version string: ", modelVisit.OsVersionString,
		)

		log.Println(
			"FirstPageTitle: ", modelVisit.FirstPageTitle,
			"FirstPageDomain: ", modelVisit.FirstPageDomain,
			"FirstPageURL: ", modelVisit.FirstPageURL,
		)

		log.Println(
			"FirstPageSourceURL: ", modelVisit.FirstPageSourceURL,
			"FirstPageSourceKeyword: ", modelVisit.FirstPageSourceKeyword,
			"FirstPageSource: ", modelVisit.FirstPageSource,
		)

		log.Printf(
			"ip: %s\n Country: %s\n Province: %s\n City: %s\n Isp: %s\n",
			modelVisit.IP,
			modelVisit.Country,
			modelVisit.Province,
			modelVisit.City,
			modelVisit.Isp,
		)
	}
}

func TestIMService_InitVisitHasTrace(t *testing.T) {
	test.InitTest()
	defer test.Clear()

	ast := assert.New(t)
	handler := &IMService{loc: test.IpLocation}
	entID := common.GenUniqueID()
	data := &InitVisitReq{
		EntID:   entID,
		TraceID: "",
		Keyword: "test key word",
		Title:   "test title",
	}

	v, err := common.Marshal(data)
	ast.Nil(err)

	_, rec, echoContext := setHttpRecorder(http.MethodPost, v)
	echoContext.SetPath("/api/v1/enterprises/:ent_id/visits")
	echoContext.SetParamNames("ent_id")
	echoContext.SetParamValues(data.EntID)

	var createVisitResp *InitVisitResp
	if ast.NoError(handler.InitVisit(echoContext)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		body := rec.Body.String()

		var s *initVisitResp
		err = common.Unmarshal(body, &s)
		ast.Nil(err)

		ast.Equal(s.Code, 0)
		createVisitResp = s.Body
	}

	data = &InitVisitReq{
		EntID:   entID,
		TraceID: createVisitResp.TraceID,
		Keyword: "new key word",
		Title:   "new title",
	}
	dataStr, err := common.Marshal(data)
	ast.Nil(err)
	_, rec, echoContext = setHttpRecorder(http.MethodPost, dataStr)
	echoContext.SetPath("/api/v1/enterprises/:ent_id/visits")
	echoContext.SetParamNames("ent_id")
	echoContext.SetParamValues(entID)
	if ast.Nil(handler.InitVisit(echoContext)) {
		body := rec.Body.String()

		var s *initVisitResp
		err = common.Unmarshal(body, &s)
		ast.Nil(err)
		ast.Equal(s.Code, 0)
		modelVisit, err := models.VisitByID(db.Mysql, s.Body.VisitID)
		ast.Nil(err)

		modelVisits, err := models.VisitsByEntIDTraceID(db.Mysql, entID, s.Body.TraceID)
		ast.Nil(err)

		ast.Equal(modelVisit, modelVisits[0])
		ast.Equal(2, modelVisit.VisitPageCnt)

		visitor, err := models.VisitorByID(db.Mysql, s.Body.VisitorID)
		ast.Nil(err)
		ast.Equal(entID, visitor.EntID)
		ast.Equal(s.Body.TraceID, visitor.TraceID)
		ast.Equal(2, visitor.VisitCnt)
		ast.Equal(2, visitor.VisitPageCnt)
	}
}
