package models

import "time"

func GetAgentStats(db XODB, entID string, start, end time.Time) (result map[int32][]*AgentStatistic, err error) {
	stats, err := AgentStatisticsByEntID(db, entID, start, end)
	if err != nil {
		return nil, err
	}

	result = map[int32][]*AgentStatistic{}
	for _, s := range stats {
		dt := getHourDateFromTime(s.CreatedAt)
		if v, ok := result[dt]; ok {
			v = append(v, s)
			result[dt] = v
			continue
		}

		result[dt] = []*AgentStatistic{s}
	}

	return
}

func getHourDateFromTime(t time.Time) int32 {
	year, month, day := t.Date()
	monthInt := int(month)
	hour := t.Hour()

	return int32(year*1000000 + monthInt*10000 + day*100 + hour)
}
