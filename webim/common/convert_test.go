package common

import (
	"fmt"
	"testing"
	"time"

	assert2 "github.com/stretchr/testify/assert"
)

func TestConvertStringSliceToInterface(t *testing.T) {
	assert := assert2.New(t)

	res := ConvertStringSliceToInterface([]string{"1", "2", "3"})
	assert.Equal([]interface{}{"1", "2", "3"}, res)
}

func TestConvertUTCToLocal(t *testing.T) {
	tm := time.Now().UTC()
	fmt.Println("tm", tm)

	x := ConvertUTCToLocal(tm)
	fmt.Println(x)
}

func TestConvertStringToTime(t *testing.T) {
	v := `2019-05-02T11:16:10.749119`
	tm, err := time.Parse(layout1, v)
	fmt.Println(tm, err)
}

func TestConvertStringToTime1(t *testing.T) {
	v := `2019-05-02`
	tm, err := time.Parse("2006-01-02", v)
	fmt.Println(tm, err)
}

func TestConvertUTCToTimeString(t *testing.T) {
	tm, _ := time.Parse(layout1, "2019-06-30T12:47:11.461550")
	ss := ConvertUTCToTimeString(tm)
	fmt.Println(*ss)
}
