package handler

import (
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"

	"bitbucket.org/forfd/custm-chat/webim/db"
	"bitbucket.org/forfd/custm-chat/webim/middleware"
	"bitbucket.org/forfd/custm-chat/webim/models"
)

type convEveryHourCount struct {
	ConvCnt      int `json:"conv_cnt"`
	MsgCnt       int `json:"msg_cnt"`
	VisitCnt     int `json:"visit_cnt"`
	VisitPageCnt int `json:"visit_page_cnt"`
	VisitorCnt   int `json:"visitor_cnt"`
}

type convEverydayHourCounts map[string]*convEveryHourCount

type everyDayCount struct {
	Date    string              `json:"date"`
	Summary *convEveryHourCount `json:"summary"`
	convEverydayHourCounts
}

// // {"begin":"2019-04-23 16","end":"2019-04-30 15","browser_id":"agent1556703681959"}
type TimeRange struct {
	Begin string `json:"begin"`
	End   string `json:"end"`
}

type ConversationStatsResp struct {
	Stats []*everyDayCount `json:"stats"`
}

type agentBrief struct {
	ID       string `json:"id"`
	Avatar   string `json:"avatar"`
	Nickname string `json:"nickname"`
	Realname string `json:"realname"`
}

type agentEvalCnt struct {
	Agent     *agentBrief `json:"agent"`
	BadCnt    int         `json:"bad_cnt"`
	ConvCnt   int         `json:"conv_cnt"`
	GoodCnt   int         `json:"good_cnt"`
	MediumCnt int         `json:"medium_cnt"`
}

type hourDate int32 // 2019050100
type everyDayEvalCnt map[hourDate][]*agentEvalCnt

type everyDayEvalCnts struct {
	day             int32
	everyDayEvalCnt everyDayEvalCnt
}

type EvaluationStatsResp struct {
	Stats []everyDayEvalCnt `json:"stats"`
}

type agentWorkload struct {
	AgentID          string `json:"agent_id"`
	AssignedConvsCnt int    `json:"assigned_convs_cnt"`
	AvaChatCnt       int    `json:"ava_chat_cnt"`
	AvgDurationTime  int    `json:"avg_duration_time"`
	AvgWaitTime      int    `json:"avg_wait_time"`
	BadConvCnt       int    `json:"bad_conv_cnt"`
	BronzeConvCnt    int    `json:"bronze_conv_cnt"`
	CluesMsg         int    `json:"clues_msg"`
	ConvCnt          int    `json:"conv_cnt"`
	DelayChatCnt     int    `json:"delay_chat_cnt"`
	DurationTime     int    `json:"duration_time"`
	EvaConvCnt       int    `json:"eva_conv_cnt"`
	GoldConvCnt      int    `json:"gold_conv_cnt"`
	GoodConvCnt      int    `json:"good_conv_cnt"`
	MediumConvCnt    int    `json:"medium_conv_cnt"`
	MissChatCnt      int    `json:"miss_chat_cnt"`
	MsgCnt           int    `json:"msg_cnt"`
	NewRemarkCnt     int    `json:"new_remark_cnt"`
	NoEvaConvCnt     int    `json:"no_eva_conv_cnt"`
	NogradeConvCnt   int    `json:"nograde_conv_cnt"`
	RedirectCnt      int    `json:"redirect_cnt"`
	RedirectedCnt    int    `json:"redirected_cnt"`
	SilverConvCnt    int    `json:"silver_conv_cnt"`
	WaitTime         int    `json:"wait_time"`
}

type agentWorkLoadStats struct {
	Stats []*agentWorkload `json:"stats"`
}

// POST /api/stats/conv
func (s *IMService) ConversationStats(ctx echo.Context) error {
	req := &TimeRange{}
	if err := ctx.Bind(req); err != nil {
		return invalidParameterResp(ctx, err.Error())
	}

	if req.Begin == "" || req.End == "" {
		return invalidParameterResp(ctx, "begin/end time invalid")
	}

	begin, end, err := parseTimeRange(req.Begin, req.End)
	if err != nil {
		return invalidParameterResp(ctx, fmt.Sprintf("invalid date format: %v", err))
	}

	entID := ctx.Get(middleware.AgentEntIDKey).(string)
	counts := map[hourDate]*everyHourCnt{}

	stats, err := models.ConversationStatByEntID(db.Mysql, entID, begin, end)
	if err != nil {
		return dbErrResp(ctx, err.Error())
	}

	for dt, v := range stats {
		counts[hourDate(dt)] = convertModelConvStatsToEveryHourCnt(v)
	}

	setDayCnt := func(v *everyDayCount, hourDt hourDate, cnt *everyHourCnt) {
		v.Summary.ConvCnt += cnt.ConvCnt
		v.Summary.MsgCnt += cnt.MsgCnt
		v.Summary.VisitCnt += cnt.VisitCnt
		v.Summary.VisitorCnt += cnt.VisitorCnt
		v.Summary.VisitPageCnt += cnt.VisitPageCnt
		v.convEverydayHourCounts[strconv.FormatInt(int64(hourDt), 10)] = &convEveryHourCount{
			ConvCnt:      cnt.ConvCnt,
			MsgCnt:       cnt.MsgCnt,
			VisitCnt:     cnt.VisitCnt,
			VisitPageCnt: cnt.VisitPageCnt,
			VisitorCnt:   cnt.VisitorCnt,
		}
	}
	resultCnt := map[int32]*everyDayCount{}
	for hourDt, cnt := range counts {
		day := int32(hourDt / 100)
		v, ok := resultCnt[day]
		if ok {
			setDayCnt(v, hourDt, cnt)
			resultCnt[day] = v
			continue
		}

		v = &everyDayCount{
			Date:                   strconv.FormatInt(int64(day), 10),
			Summary:                &convEveryHourCount{},
			convEverydayHourCounts: map[string]*convEveryHourCount{},
		}
		setDayCnt(v, hourDt, cnt)
		resultCnt[day] = v
	}

	var everyDayCnt = make([]*everyDayCount, 0, len(resultCnt))
	for _, cnt := range resultCnt {
		everyDayCnt = append(everyDayCnt, cnt)
	}

	sort.SliceStable(everyDayCnt, func(i, j int) bool {
		return everyDayCnt[i].Date <= everyDayCnt[j].Date
	})

	result := &ConversationStatsResp{Stats: everyDayCnt}
	return jsonResponse(ctx, result)
}

// POST /api/stats/evaluation
func (s *IMService) EvaluationStats(ctx echo.Context) error {
	req := &TimeRange{}
	if err := ctx.Bind(req); err != nil {
		return invalidParameterResp(ctx, err.Error())
	}

	if req.Begin == "" || req.End == "" {
		return invalidParameterResp(ctx, "begin/end time invalid")
	}

	start, end, err := parseTimeRange(req.Begin, req.End)
	if err != nil {
		return invalidParameterResp(ctx, err.Error())
	}

	entID := ctx.Get(middleware.AgentEntIDKey).(string)
	agents, err := models.AgentsByEntID(db.Mysql, entID)
	if err != nil {
		return dbErrResp(ctx, err.Error())
	}

	briefAgents := map[string]*agentBrief{}
	for _, agent := range agents {
		briefAgents[agent.ID] = &agentBrief{ID: agent.ID, Avatar: agent.Avatar, Nickname: agent.NickName, Realname: agent.RealName}
	}

	agentStats, err := models.GetAgentStats(db.Mysql, entID, start, end)
	if err != nil {
		return dbErrResp(ctx, err.Error())
	}

	counts := map[hourDate][]*agentEvalCnt{}
	for dt, stats := range agentStats {
		hourDate := hourDate(dt)

		var hourCounts []*agentEvalCnt
		for _, st := range stats {
			hourCounts = append(hourCounts, &agentEvalCnt{
				Agent:     briefAgents[st.AgentID],
				ConvCnt:   int(st.ConversationCount),
				GoodCnt:   int(st.GoodCount),
				MediumCnt: int(st.MediumCount),
				BadCnt:    int(st.BadCount),
			})
		}
		counts[hourDate] = hourCounts
	}

	var resultCnt = map[int32]everyDayEvalCnt{}
	for hourDt, cnt := range counts {
		day := int32(hourDt / 100)
		v, ok := resultCnt[day]
		if ok {
			v[hourDt] = cnt
			resultCnt[day] = v
			continue
		}

		resultCnt[day] = everyDayEvalCnt{hourDt: cnt}
	}

	result := &EvaluationStatsResp{Stats: []everyDayEvalCnt{}}
	var everyDayCnt = make([]everyDayEvalCnts, 0, len(resultCnt))
	for day, cnt := range resultCnt {
		everyDayCnt = append(everyDayCnt, everyDayEvalCnts{day: day, everyDayEvalCnt: cnt})
	}

	sort.SliceStable(everyDayCnt, func(i, j int) bool {
		return everyDayCnt[i].day <= everyDayCnt[j].day
	})
	for _, cnt := range everyDayCnt {
		result.Stats = append(result.Stats, cnt.everyDayEvalCnt)
	}

	return jsonResponse(ctx, result)
}

// POST /api/analytics/agent/workload
func (s *IMService) AgentWorkLoad(ctx echo.Context) error {
	req := &TimeRangeInTimestamp{}
	if err := ctx.Bind(req); err != nil {
		return invalidParameterResp(ctx, err.Error())
	}

	entID := ctx.Get(middleware.AgentEntIDKey).(string)
	agentIDs, err := models.AgentIDsByEntID(db.Mysql, entID)
	if err != nil {
		return dbErrResp(ctx, err.Error())
	}

	start, end := time.Unix(req.Begin/1000, 0), time.Unix(req.End/1000, 0)
	evalStats, err := models.EvalConversationStatsByTimeRange(db.Mysql, agentIDs, start, end)
	if err != nil {
		return dbErrResp(ctx, err.Error())
	}

	qualityStats, err := models.QualityConversationStatsByTimeRange(db.Mysql, agentIDs, start, end)
	if err != nil {
		return dbErrResp(ctx, err.Error())
	}

	overallConvStats, err := models.OverallConversationStatsByTimeRange(db.Mysql, entID, start, end)
	if err != nil {
		return dbErrResp(ctx, err.Error())
	}

	result := &agentWorkLoadStats{Stats: []*agentWorkload{}}
	for _, id := range agentIDs {
		workLoad := &agentWorkload{AgentID: id}
		evalStat, ok := evalStats[id]
		if ok {
			workLoad.EvaConvCnt = evalStat.EvaConvCnt
			workLoad.GoodConvCnt = evalStat.GoodConvCnt
			workLoad.MediumConvCnt = evalStat.MediumConvCnt
			workLoad.BadConvCnt = evalStat.BadConvCnt
			workLoad.NoEvaConvCnt = evalStat.NoEvaConvCnt
		}

		qualityStat, ok := qualityStats[id]
		if ok {
			workLoad.GoldConvCnt = qualityStat.GoldConvCnt
			workLoad.SilverConvCnt = qualityStat.SilverConvCnt
			workLoad.BronzeConvCnt = qualityStat.BronzeConvCnt
			workLoad.NogradeConvCnt = qualityStat.NogradeConvCnt
		}

		overAllStat, ok := overallConvStats[id]
		if ok {
			workLoad.ConvCnt = overAllStat.ConvCount
			workLoad.AvgDurationTime = overAllStat.AvgDurationInSec
			workLoad.DurationTime = overAllStat.DurationInSec
			workLoad.MsgCnt = overAllStat.MsgCnt
		}

		result.Stats = append(result.Stats, workLoad)
	}

	return jsonResponse(ctx, result)
}
