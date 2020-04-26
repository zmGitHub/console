package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestGetOffsetLimit(t *testing.T) {
	ast := assert.New(t)
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/?offset=1&limit=100", nil)
	rec := httptest.NewRecorder()
	echoContext := e.NewContext(req, rec)
	offset, limit, err := getOffsetLimitFromCtx(echoContext)
	if ast.Nil(err) {
		ast.Equal(1, offset)
		ast.Equal(100, limit)
	}

	req = httptest.NewRequest(http.MethodGet, "/?offset=1&limit=", nil)
	rec = httptest.NewRecorder()
	echoContext = e.NewContext(req, rec)
	offset, limit, err = getOffsetLimitFromCtx(echoContext)
	if ast.Nil(err) {
		ast.Equal(1, offset)
		ast.Equal(30, limit)
	}

	req = httptest.NewRequest(http.MethodGet, "/?offset=&limit=10", nil)
	rec = httptest.NewRecorder()
	echoContext = e.NewContext(req, rec)
	offset, limit, err = getOffsetLimitFromCtx(echoContext)
	if ast.Nil(err) {
		ast.Equal(0, offset)
		ast.Equal(10, limit)
	}
}

type rankTester []int

func (r rankTester) Rank(idx int) int {
	return r[idx]
}

func (r rankTester) Length() int {
	return len(r)
}

func TestGetNewRank(t *testing.T) {
	r := rankTester([]int{100000, 200000, 300000})
	newRank := getNewRank(0, 1, r)
	assert.Equal(t, 250000, newRank)

	newRank = getNewRank(0, 0, r)
	assert.Equal(t, -1, newRank)

	newRank = getNewRank(0, 2, r)
	assert.Equal(t, 400000, newRank)

	newRank = getNewRank(0, 2, r)
	assert.Equal(t, 400000, newRank)

	newRank = getNewRank(1, 1, r)
	assert.Equal(t, -1, newRank)

	newRank = getNewRank(1, 2, r)
	assert.Equal(t, 400000, newRank)

	newRank = getNewRank(2, 0, r)
	assert.Equal(t, 50000, newRank)

	newRank = getNewRank(1, 0, r)
	assert.Equal(t, 50000, newRank)

	newRank = getNewRank(2, 1, r)
	assert.Equal(t, 150000, newRank)
}
