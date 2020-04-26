package handler

import (
	"sort"
	"time"

	"github.com/labstack/echo/v4"

	"bitbucket.org/forfd/custm-chat/webim/db"
	"bitbucket.org/forfd/custm-chat/webim/middleware"
	"bitbucket.org/forfd/custm-chat/webim/models"
)

var dayInterval = 5 // 5 days
var defaultTrafficCnt = func() *trafficCnt {
	return &trafficCnt{DeskDirect: &visitConvCnt{}}
}

type TimeRangeInTimestamp struct {
	Begin int64 `json:"begin"`
	End   int64 `json:"end"`
}

type tmRange struct {
	start time.Time
	end   time.Time
}

type tmRanges []*tmRange

type everyHourCnt struct {
	AcceptInviteCnt  int `json:"accept_invite_cnt"`
	ActiveConvCnt    int `json:"active_conv_cnt"`
	BadConvCnt       int `json:"bad_conv_cnt"`
	BronzeConvCnt    int `json:"bronze_conv_cnt"`
	CluesCnt         int `json:"clues_cnt"`
	ConvCnt          int `json:"conv_cnt"`
	ConvDuration     int `json:"conv_duration"`
	Duration         int `json:"duration"`
	EffectiveConvCnt int `json:"effective_conv_cnt"`
	GoldConvCnt      int `json:"gold_conv_cnt"`
	GoodConvCnt      int `json:"good_conv_cnt"`
	InvitationCnt    int `json:"invitation_cnt"`
	MediumConvCnt    int `json:"medium_conv_cnt"`
	MsgCnt           int `json:"msg_cnt"`
	NogradeConvCnt   int `json:"nograde_conv_cnt"`
	RejectInviteCnt  int `json:"reject_invite_cnt"`
	RemarksCnt       int `json:"remarks_cnt"`
	SilverConvCnt    int `json:"silver_conv_cnt"`
	VisitCnt         int `json:"visit_cnt"`
	VisitPageCnt     int `json:"visit_page_cnt"`
	VisitorCnt       int `json:"visitor_cnt"`
	WaitTime         int `json:"wait_time"`
}

type OverViewResponse struct {
	Stats map[hourDate]*everyHourCnt `json:"stats"`
}

type visitConvCnt struct {
	ConvCnt            int `json:"conv_cnt"`
	IneffectiveConvCnt int `json:"ineffective_conv_cnt"`
	MsgCnt             int `json:"msg_cnt"`
	VisitCnt           int `json:"visit_cnt"`
	VisitPageCnt       int `json:"visit_page_cnt"`
	VisitorCnt         int `json:"visitor_cnt"`
}

type trafficCnt struct {
	Appsdk         *visitConvCnt `json:"appsdk"`
	Campaign       *visitConvCnt `json:"campaign"`
	DeskDirect     *visitConvCnt `json:"desk_direct"`
	DeskReferral   *visitConvCnt `json:"desk_referral"`
	DeskSearch     *visitConvCnt `json:"desk_search"`
	MobileDirect   *visitConvCnt `json:"mobile_direct"`
	MobileReferral *visitConvCnt `json:"mobile_referral"`
	MobileSearch   *visitConvCnt `json:"mobile_search"`
	Toutiao        *visitConvCnt `json:"toutiao"`
	Weibo          *visitConvCnt `json:"weibo"`
	Weixin         *visitConvCnt `json:"weixin"`
}

type everyHourTrafficCnt map[hourDate]*trafficCnt

type everyDayTrafficCnt struct {
	day                 int32
	everyHourTrafficCnt everyHourTrafficCnt
}

type TrafficOverviewResp struct {
	Stats []everyHourTrafficCnt `json:"stats"`
}

// POST /api/analytics/overview
func (s *IMService) OverView(ctx echo.Context) (err error) {
	req := &TimeRangeInTimestamp{}
	if err = ctx.Bind(req); err != nil {
		return invalidParameterResp(ctx, err.Error())
	}

	if req.Begin <= 0 || req.End <= 0 {
		return invalidParameterResp(ctx, "invalid begin/end time stamp")
	}

	start, end := time.Unix(req.Begin/1000, 0), time.Unix(req.End/1000, 0)

	entID := ctx.Get(middleware.AgentEntIDKey).(string)

	resp := &OverViewResponse{}
	resp.Stats, err = s.convVisitOverview(entID, start, end)
	if err != nil {
		return internalServerErr(ctx, err.Error())
	}

	return jsonResponse(ctx, resp)
}

// POST /api/stats/traffic_overview
// {"begin":"2019-04-30 16","end":"2019-05-01 15","browser_id":"agent1556712701751"}
func (s *IMService) TrafficOverview(ctx echo.Context) error {
	req := &TimeRange{}
	if err := ctx.Bind(req); err != nil {
		return invalidParameterResp(ctx, err.Error())
	}

	begin, end, err := parseTimeRange(req.Begin, req.End)
	if err != nil {
		return invalidParameterResp(ctx, err.Error())
	}

	entID := ctx.Get(middleware.AgentEntIDKey).(string)
	counts := everyHourTrafficCnt{}
	tmRanges := toTimeRanges(begin, end)
	for _, tmr := range tmRanges {
		convCounts, visitCounts, visitorCounts, pageCounts, err := getConvVisitRecords(entID, tmr.start, tmr.end)
		if err != nil {
			return dbErrResp(ctx, err.Error())
		}

		setEveryHourTrafficCount(counts, convCounts, visitCounts, visitorCounts, pageCounts)
	}

	// map[day]everyHourTrafficCnt
	var resultCnt = map[int32]everyHourTrafficCnt{}
	for hourDt, cnt := range counts {
		day := int32(hourDt / 100)
		v, ok := resultCnt[day]
		if ok {
			v[hourDt] = cnt
			resultCnt[day] = v
			continue
		}

		resultCnt[day] = everyHourTrafficCnt{hourDt: cnt}
	}

	var everyDayCnt = make([]*everyDayTrafficCnt, 0, len(resultCnt))
	for day, cnt := range resultCnt {
		everyDayCnt = append(everyDayCnt, &everyDayTrafficCnt{day: day, everyHourTrafficCnt: cnt})
	}

	sort.SliceStable(everyDayCnt, func(i, j int) bool {
		return everyDayCnt[i].day <= everyDayCnt[j].day
	})
	result := &TrafficOverviewResp{Stats: []everyHourTrafficCnt{}}
	for _, cnt := range everyDayCnt {
		result.Stats = append(result.Stats, cnt.everyHourTrafficCnt)
	}

	return jsonResponse(ctx, result)
}

func getConvVisitRecords(entID string, start, end time.Time) (conversations []*models.Conversation, visits []*models.Visit, visitors []*models.Visitor, pages []*models.VisitPage, err error) {
	conversations, err = models.ConversationsByTimeRange(db.Mysql, entID, start, end)
	if err != nil {
		return
	}

	visits, err = models.VisitsByTimeRange(db.Mysql, entID, start, end)
	if err != nil {
		return
	}

	visitors, err = models.VisitorByTimeRange(db.Mysql, entID, start, end)
	if err != nil {
		return
	}

	pages, err = models.VisitPageByTimeRange(db.Mysql, entID, start, end)
	if err != nil {
		return
	}

	return
}

func (s *IMService) convVisitOverview(entID string, start, end time.Time) (map[hourDate]*everyHourCnt, error) {
	counts := map[hourDate]*everyHourCnt{}

	//tmRanges := toTimeRanges(start, end)
	//for _, tmr := range tmRanges {
	//	conversations, visits, visitors, pages, err := getConvVisitRecords(entID, tmr.start, tmr.end)
	//	if err != nil {
	//		return nil, err
	//	}
	//
	//	setEveryHourConvCount(counts, conversations)
	//	setEveryHourVisitCount(counts, visits, visitors, pages)
	//}
	result, err := models.ConversationStatByEntID(db.Mysql, entID, start, end)
	if err != nil {
		return nil, err
	}

	for dt, v := range result {
		counts[hourDate(dt)] = convertModelConvStatsToEveryHourCnt(v)
	}

	return counts, nil
}

func getHourDate(t time.Time) hourDate {
	year, month, day := t.Date()
	monthInt := int(month)
	hour := t.Hour()

	return hourDate(year*1000000 + monthInt*10000 + day*100 + hour)
}

func setEveryHourConvCount(counts map[hourDate]*everyHourCnt, conversations []*models.Conversation) {
	setCnt := func(v *everyHourCnt, conv *models.Conversation) *everyHourCnt {
		v.AcceptInviteCnt += 1
		v.ActiveConvCnt += 1
		v.ConvCnt += 1
		v.EffectiveConvCnt += 1
		v.ConvDuration += int(conv.Duration)
		v.Duration += int(conv.Duration)

		switch conv.EvalLevel {
		case models.GoodConversation:
			v.GoodConvCnt += 1
		case models.MediumConversation:
			v.MediumConvCnt += 1
		case models.BadConversation:
			v.BadConvCnt += 1
		}

		switch conv.QualityGrade.String {
		case models.FirstLevel:
			v.GoldConvCnt += 1
		case models.SecondLevel:
			v.SilverConvCnt += 1
		case models.ThirdLevel:
			v.BronzeConvCnt += 1
		default:
			v.NogradeConvCnt += 1
		}

		v.MsgCnt += int(conv.MsgCount)
		return v
	}

	for _, conv := range conversations {
		hourDate := getHourDate(conv.CreatedAt)
		if v, ok := counts[hourDate]; ok {
			counts[hourDate] = setCnt(v, conv)
			continue
		}

		hourCount := &everyHourCnt{}
		counts[hourDate] = setCnt(hourCount, conv)
	}
}

func setEveryHourTrafficCount(cnt everyHourTrafficCnt, conversations []*models.Conversation, visits []*models.Visit, visitors []*models.Visitor, pages []*models.VisitPage) {
	setCnt := func(v *trafficCnt, conv *models.Conversation) {
		v.DeskDirect.ConvCnt += 1
		v.DeskDirect.MsgCnt += int(conv.MsgCount)
	}
	for _, conv := range conversations {
		dt := getHourDate(conv.CreatedAt)
		if v, ok := cnt[dt]; ok {
			setCnt(v, conv)
			continue
		}

		hourCount := defaultTrafficCnt()
		setCnt(hourCount, conv)
		cnt[dt] = hourCount
	}

	for _, visit := range visits {
		dt := getHourDate(visit.CreatedAt)
		v, ok := cnt[dt]
		if ok {
			v.DeskDirect.VisitCnt += 1
			continue
		}

		hourCount := defaultTrafficCnt()
		hourCount.DeskDirect.VisitCnt += 1
		cnt[dt] = hourCount
	}

	for _, visitor := range visitors {
		dt := getHourDate(visitor.CreatedAt)
		v, ok := cnt[dt]
		if ok {
			v.DeskDirect.VisitorCnt += 1
			continue
		}

		hourCount := defaultTrafficCnt()
		hourCount.DeskDirect.VisitorCnt += 1
		cnt[dt] = hourCount
	}

	for _, page := range pages {
		dt := getHourDate(page.CreatedAt)
		v, ok := cnt[dt]
		if ok {
			v.DeskDirect.VisitPageCnt += 1
			continue
		}

		hourCount := defaultTrafficCnt()
		hourCount.DeskDirect.VisitPageCnt += 1
		cnt[dt] = hourCount
	}
}

func setEveryHourVisitCount(counts map[hourDate]*everyHourCnt, visits []*models.Visit, visitors []*models.Visitor, pages []*models.VisitPage) {
	for _, visit := range visits {
		dt := getHourDate(visit.CreatedAt)
		v, ok := counts[dt]
		if ok {
			v.VisitCnt += 1
			continue
		}

		hourCount := &everyHourCnt{}
		hourCount.VisitCnt += 1
		counts[dt] = hourCount
	}

	for _, visitor := range visitors {
		dt := getHourDate(visitor.CreatedAt)
		v, ok := counts[dt]
		if ok {
			v.VisitorCnt += 1
			if visitor.Remark != "" {
				v.RemarksCnt += 1
			}
			continue
		}

		hourCount := &everyHourCnt{}
		hourCount.VisitorCnt += 1
		if visitor.Remark != "" {
			hourCount.RemarksCnt += 1
		}
		counts[dt] = hourCount
	}

	for _, page := range pages {
		dt := getHourDate(page.CreatedAt)
		v, ok := counts[dt]
		if ok {
			v.VisitPageCnt += 1
			continue
		}

		hourCount := &everyHourCnt{}
		hourCount.VisitPageCnt += 1
		counts[dt] = hourCount
	}
}

func toTimeRanges(start, end time.Time) (ranges []*tmRange) {
	hours := end.Sub(start).Hours()
	if hours <= float64(dayInterval*24) {
		return []*tmRange{{start: start, end: end}}
	}

	count := int64(hours) / int64(dayInterval*24)
	interval := time.Duration(dayInterval) * 24 * time.Hour

	start1 := start
	end1 := start.Add(interval)
	ranges = append(ranges, &tmRange{start: start1, end: end1})

	var i int64 = 1
	for i < count {
		start1 = end1
		end1 = start1.Add(interval)
		ranges = append(ranges, &tmRange{start: start1, end: end1})
		i++
	}

	rg := &tmRange{start: end1, end: end}
	ranges = append(ranges, rg)
	return
}

func convertModelConvStatsToEveryHourCnt(v *models.ConversationStat) *everyHourCnt {
	return &everyHourCnt{
		AcceptInviteCnt:  0,
		ActiveConvCnt:    int(v.EffectiveCount),
		BadConvCnt:       int(v.BadCount),
		BronzeConvCnt:    int(v.BronzeCount),
		ConvCnt:          int(v.TotalCount),
		ConvDuration:     int(v.DurationInSec),
		Duration:         int(v.DurationInSec),
		EffectiveConvCnt: int(v.EffectiveCount),
		GoldConvCnt:      int(v.GoldCount),
		GoodConvCnt:      int(v.GoodCount),
		InvitationCnt:    0,
		MediumConvCnt:    int(v.MediumCount),
		MsgCnt:           int(v.MessageCount),
		NogradeConvCnt:   int(v.NoGradeCount),
		RemarksCnt:       int(v.RemarkCount),
		SilverConvCnt:    int(v.SilverCount),
		VisitCnt:         int(v.VisitCount),
		VisitPageCnt:     int(v.VisitPageCount),
		VisitorCnt:       int(v.VisitorCount),
	}
}
