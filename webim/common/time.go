package common

import (
	"strconv"
	"time"
)

func TimeStampToTime(stamp string) (time.Time, error) {
	i, err := strconv.ParseInt(stamp, 10, 64)
	if err != nil {
		return time.Time{}, err
	}

	tm := time.Unix(i, 0)
	return tm, nil
}
