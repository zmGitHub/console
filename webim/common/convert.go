package common

import (
	"time"

	"github.com/jinzhu/now"
)

var (
	layout     = "2006-01-02 15:04:05.999999999"
	layout1    = "2006-01-02T15:04:05.999999999"
	cnTimeZone = "Asia/Shanghai"
)

// ConvertStringSliceToInterface ...
func ConvertStringSliceToInterface(values []string) (res []interface{}) {
	for _, v := range values {
		res = append(res, v)
	}
	return
}

func ConvertUTCToTimeString(t time.Time) *string {
	if t.IsZero() {
		return nil
	}

	localTime := ConvertUTCToLocal(t)
	newTime := localTime.Format(layout)
	return &newTime
}

func ConvertUTCToLocal(t time.Time) time.Time {
	if t.IsZero() {
		return time.Time{}
	}

	location, err := time.LoadLocation(cnTimeZone)
	if err != nil {
		return t
	}

	localTime := t.In(location)
	tt, err := ConvertTimeStringToTime(localTime.Format(layout))
	if err != nil {
		return t
	}

	return tt
}

func ConvertTimeStringToTime(ts string) (t time.Time, err error) {
	t, err = now.Parse(ts)
	return
}
