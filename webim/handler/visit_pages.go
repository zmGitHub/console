package handler

import (
	"github.com/labstack/echo/v4"

	"bitbucket.org/forfd/custm-chat/webim/common"
	"bitbucket.org/forfd/custm-chat/webim/db"
	"bitbucket.org/forfd/custm-chat/webim/dto"
	"bitbucket.org/forfd/custm-chat/webim/middleware"
	"bitbucket.org/forfd/custm-chat/webim/models"
)

type VisitPageResp struct {
	Pages []struct {
		CreatedOn     string `json:"created_on"`
		Source        string `json:"source"`
		SourceKeyword string `json:"source_keyword"`
		SourceURL     string `json:"source_url"`
		Title         string `json:"title"`
		URL           string `json:"url"`
	} `json:"pages"`
}

//GET  /api/client/:track_id/pages
func (s *IMService) GetConversationPages(ctx echo.Context) (err error) {
	trackID := ctx.Param("track_id")
	visits, err := models.VisitsByTraceIDs(db.Mysql, []string{trackID})
	if err != nil {
		return dbErrResp(ctx, err.Error())
	}

	resp := &dto.VisitPages{
		Pages: []*dto.VisitPage{},
	}

	if len(visits) > 0 {
		entID := ctx.Get(middleware.AgentEntIDKey).(string)
		pages, err := models.VisitPagesByEntIDVisitID(db.Mysql, entID, visits[0].ID)
		if err != nil {
			return dbErrResp(ctx, err.Error())
		}

		for _, page := range pages {
			resp.Pages = append(resp.Pages, &dto.VisitPage{
				CreatedOn:     *common.ConvertUTCToTimeString(page.CreatedAt),
				Source:        page.Source,
				SourceKeyword: page.SourceKeyword,
				SourceURL:     page.SourceURL,
				Title:         page.Title,
				URL:           page.SourceURL,
			})
		}
	}

	return jsonResponse(ctx, resp)
}

// GET /api/client/1EACfgGhNoogG9YFlTVr2OJt9lK/attrs
func (s *IMService) GetAttrs(ctx echo.Context) (err error) {
	return nil
}
