package handler

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestToTimeRanges(t *testing.T) {
	cases := []struct {
		start  time.Time
		end    time.Time
		ranges []*tmRange
	}{
		{
			start: time.Date(2019, 5, 6, 8, 0, 0, 0, time.UTC),
			end:   time.Date(2019, 5, 6, 18, 0, 0, 0, time.UTC),
			ranges: []*tmRange{
				{
					start: time.Date(2019, 5, 6, 8, 0, 0, 0, time.UTC),
					end:   time.Date(2019, 5, 6, 18, 0, 0, 0, time.UTC),
				},
			},
		},
		{
			start: time.Date(2019, 5, 6, 8, 0, 0, 0, time.UTC),
			end:   time.Date(2019, 5, 12, 18, 0, 0, 0, time.UTC),
			ranges: []*tmRange{
				{
					start: time.Date(2019, 5, 6, 8, 0, 0, 0, time.UTC),
					end:   time.Date(2019, 5, 11, 8, 0, 0, 0, time.UTC),
				},
				{
					start: time.Date(2019, 5, 11, 8, 0, 0, 0, time.UTC),
					end:   time.Date(2019, 5, 12, 18, 0, 0, 0, time.UTC),
				},
			},
		},
		{
			start: time.Date(2019, 5, 5, 8, 0, 0, 0, time.UTC),
			end:   time.Date(2019, 5, 16, 18, 0, 0, 0, time.UTC),
			ranges: []*tmRange{
				{
					start: time.Date(2019, 5, 5, 8, 0, 0, 0, time.UTC),
					end:   time.Date(2019, 5, 10, 8, 0, 0, 0, time.UTC),
				},
				{
					start: time.Date(2019, 5, 10, 8, 0, 0, 0, time.UTC),
					end:   time.Date(2019, 5, 15, 8, 0, 0, 0, time.UTC),
				},
				{
					start: time.Date(2019, 5, 15, 8, 0, 0, 0, time.UTC),
					end:   time.Date(2019, 5, 16, 18, 0, 0, 0, time.UTC),
				},
			},
		},
	}

	for _, c := range cases {
		ranges := toTimeRanges(c.start, c.end)

		for _, tm := range ranges {
			fmt.Println(tm.start, tm.end)
		}

		assert.Equal(t, c.ranges, ranges)
	}
}

func TestToTimeRanges1(t *testing.T) {
	startStamp := int64(1554566400000 / 1000)
	endStamp := int64(1557244799000 / 1000)

	fmt.Println(1556553600000, time.Now().Unix())

	start, end := time.Unix(startStamp, 0), time.Unix(endStamp, 0)
	fmt.Println(start, end)

	ranges := toTimeRanges(start, end)
	// fmt.Println(ranges)

	for _, rg := range ranges {
		fmt.Println(rg.start, rg.end)
	}
}
