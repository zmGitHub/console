package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"

	"bitbucket.org/forfd/custm-chat/webim/common"
	"bitbucket.org/forfd/custm-chat/webim/middleware"
	"bitbucket.org/forfd/custm-chat/webim/models"
)

var (
	defaultOffset = 0
	defaultLimit  = 30
)

var hasPerm = func(entID, agentID, appName, permName string) *ErrMsg {
	//allow, err := models.GetAgentPermFromCache(db.Mysql, entID, agentID, appName, permName)
	//if err != nil {
	//	return &ErrMsg{Code: common.DBErr, Message: err.Error()}
	//}
	//
	//if !allow {
	//	return &ErrMsg{Code: common.PermissionLimited, Message: "NO " + permName + " Permission"}
	//}

	return nil
}

type AgentInfo struct {
	EntID  string
	UserID string
	Token  string
}

var addAllowOriginHeader = func(ctx echo.Context) echo.Context {
	resp := ctx.Response()
	resp.Header().Set(echo.HeaderAccessControlAllowOrigin, "*")
	resp.Header().Set(echo.HeaderAccessControlAllowCredentials, "true")
	methods := []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete}
	resp.Header().Set(echo.HeaderAccessControlRequestMethod, strings.Join(methods, ", "))
	return ctx
}

var jsonResponse = func(ctx echo.Context, msg interface{}) error {
	return ctx.JSON(http.StatusOK, msg)
}

var invalidParameterResp = func(ctx echo.Context, msg string) error {
	ctx = addAllowOriginHeader(ctx)
	return ctx.JSON(http.StatusBadRequest, &ErrMsg{Code: common.InvalidParameterErr, Message: msg})
}

var dbErrResp = func(ctx echo.Context, msg string) error {
	ctx = addAllowOriginHeader(ctx)
	return jsonResponse(ctx, &ErrMsg{Code: common.DBErr, Message: msg})
}

var internalServerErr = func(ctx echo.Context, msg string) error {
	ctx = addAllowOriginHeader(ctx)
	return ctx.JSON(http.StatusInternalServerError, &ErrMsg{Code: common.InternalServerErr, Message: msg})
}

var errResp = func(ctx echo.Context, code int, msg string) error {
	ctx = addAllowOriginHeader(ctx)
	return jsonResponse(ctx, &ErrMsg{Code: code, Message: msg})
}

var noPermResp = func(ctx echo.Context, msg *ErrMsg) error {
	ctx = addAllowOriginHeader(ctx)
	return ctx.JSON(http.StatusForbidden, msg)
}

type Ranker interface {
	Rank(index int) int
	Length() int
}

type QuickReplyGroupsRankImpl []*models.QuickreplyGroup

type ClientTagsRankImpl []*models.VisitorTag

func (groupsRank QuickReplyGroupsRankImpl) Rank(index int) int {
	return groupsRank[index].Rank
}

func (groupsRank QuickReplyGroupsRankImpl) Length() int {
	return len(groupsRank)
}

func (tagRank ClientTagsRankImpl) Rank(index int) int {
	return tagRank[index].Rank
}

func (tagRank ClientTagsRankImpl) Length() int {
	return len(tagRank)
}

func getAgentInfoFromJwtToken(ctx echo.Context) *AgentInfo {
	return &AgentInfo{
		EntID:  ctx.Get(middleware.AgentEntIDKey).(string),
		UserID: ctx.Get(middleware.AgentIDKey).(string),
		Token:  ctx.Get(middleware.AgentTokenKey).(string),
	}
}

func getOffsetLimitFromCtx(ctx echo.Context) (int, int, error) {
	offset := ctx.QueryParam("offset")
	limit := ctx.QueryParam("limit")
	offsetVal, limitVal := defaultOffset, defaultLimit
	if offset != "" {
		v, err := strconv.Atoi(offset)
		if err != nil {
			return -1, -1, err
		}

		if v >= 0 {
			offsetVal = v
		}
	}

	if limit != "" {
		v, err := strconv.Atoi(limit)
		if err != nil {
			return -1, -1, err
		}

		if v >= 1 {
			limitVal = v
		}
	}

	return offsetVal, limitVal, nil
}

func getNewRank(currentPosition, targetPosition int, ranker Ranker) int {
	if currentPosition == targetPosition || currentPosition < 0 {
		return -1
	}

	if ranker.Length() > 1 {
		p := targetPosition

		if p >= 0 {
			if p == 0 {
				return ranker.Rank(0) / 2
			}

			if p < ranker.Length()-1 {
				if p > currentPosition {
					return (ranker.Rank(p) + ranker.Rank(p+1)) / 2
				}
				return (ranker.Rank(p) + ranker.Rank(p-1)) / 2
			}
			return ranker.Rank(ranker.Length()-1) + 100000
		}
	}

	return -1
}

func parseDate(s string) (time.Time, error) {
	t := time.Time{}
	ds := strings.Split(strings.Trim(s, " "), " ")
	if len(ds) != 2 {
		return t, fmt.Errorf("wrong date format")
	}

	dt := strings.Split(ds[0], "-")
	if len(dt) != 3 {
		return t, fmt.Errorf("wrong date format")
	}

	year, month, day := dt[0], dt[1], dt[2]
	yearInt, err := strconv.Atoi(year)
	if err != nil {
		return t, err
	}

	monthInt, err := strconv.Atoi(month)
	if err != nil {
		return t, err
	}

	dayInt, err := strconv.Atoi(day)
	if err != nil {
		return t, err
	}

	hour, err := strconv.Atoi(ds[1])
	if err != nil {
		return t, err
	}

	return time.Date(yearInt, time.Month(monthInt), dayInt, hour, 0, 0, 0, time.UTC), nil
}

func parseTimeRange(begin, end string) (beginTm, endTm time.Time, err error) {
	beginTm, err = parseDate(begin)
	if err != nil {
		return
	}

	endTm, err = parseDate(end)
	return
}
