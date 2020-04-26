package handler

import (
	"bytes"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/dchest/captcha"
	"github.com/labstack/echo/v4"

	"bitbucket.org/forfd/custm-chat/webim/common"
	"bitbucket.org/forfd/custm-chat/webim/common/logging"
	"bitbucket.org/forfd/custm-chat/webim/conf"
	"bitbucket.org/forfd/custm-chat/webim/db"
	"bitbucket.org/forfd/custm-chat/webim/external/submail"
	"bitbucket.org/forfd/custm-chat/webim/models"
)

const resetPasswordTemp = `%s/reset-password?token=%s&email=%s&ent_id=%s`

type RequestResetPwdReq struct {
	Email string `json:"email"`
}

type CreateCaptchasResp struct {
	ImageURL     string `json:"captcha_image_url"`
	CaptchaToken string `json:"captcha_token"`
}

type verifyCaptchaErrMsg struct {
	*ErrMsg
	CaptchaNeeded bool `json:"captcha_needed"`
}

// POST /api/captchas
func (s *IMService) CreateCaptchas(ctx echo.Context) error {
	token := captcha.NewLen(6)
	return jsonResponse(ctx, &CreateCaptchasResp{
		ImageURL:     fmt.Sprintf("/captcha_images/%s", token),
		CaptchaToken: token,
	})
}

// GET /captcha_images/:captcha_id
func (s *IMService) GetCaptcha(ctx echo.Context) error {
	id := ctx.Param("captcha_id")

	w := ctx.Response()
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	w.Header().Set("Content-Type", "image/png")

	var content bytes.Buffer
	if err := captcha.WriteImage(&content, id, captcha.StdWidth, captcha.StdHeight); err != nil {
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	r := ctx.Request()
	http.ServeContent(w, r, id, time.Time{}, bytes.NewReader(content.Bytes()))
	return nil
}

// POST /api/request_reset
func (s *IMService) RequestResetPwd(ctx echo.Context) error {
	req := &RequestResetPwdReq{}
	if err := ctx.Bind(req); err != nil {
		return invalidParameterResp(ctx, err.Error())
	}

	if req.Email == "" {
		return invalidParameterResp(ctx, "invalid email")
	}

	header := ctx.Request().Header

	token := header.Get("captcha-token")
	if token == "" {
		return invalidParameterResp(ctx, "invalid captcha token")
	}

	// captcha-value
	value := header.Get("captcha-value")
	if value == "" {
		return invalidParameterResp(ctx, "invalid captcha value")
	}

	if !captcha.VerifyString(token, value) {
		return invalidParameterResp(ctx, "图形验证码错误,请重新输入")
	}

	agent, err := models.AgentByEmail(db.Mysql, req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return invalidParameterResp(ctx, "用户不存在")
		}

		return invalidParameterResp(ctx, err.Error())
	}

	resetToken := common.RandStringBytesMask(64)
	if _, err := db.RedisClient.Set(resetToken, req.Email, 24*time.Hour).Result(); err != nil {
		return errResp(ctx, common.InternalServerErr, "服务错误，请重试")
	}

	entID := agent.EntID
	err = submail.SendEmail(
		req.Email,
		"重置chat186密码确认",
		"请打开下面的链接重置密码, 此邮件 24 小时之内有效。",
		fmt.Sprintf(resetPasswordTemp, conf.IMConf.Host, resetToken, req.Email, entID),
	)
	if err != nil {
		log.Logger.Warnf("send email error: %v", err)
	}

	return jsonResponse(ctx, &SuccessResp{Success: true})
}

func verifyCaptcha(ctx echo.Context, tokenName, tokenValueName string) *verifyCaptchaErrMsg {
	header := ctx.Request().Header
	token := header.Get(tokenName)
	if token == "" {
		return &verifyCaptchaErrMsg{ErrMsg: &ErrMsg{Message: "invalid captcha token"}, CaptchaNeeded: true}
	}

	value := header.Get(tokenValueName)
	if value == "" {
		return &verifyCaptchaErrMsg{ErrMsg: &ErrMsg{Message: "invalid captcha value"}, CaptchaNeeded: true}
	}

	if !captcha.VerifyString(token, value) {
		return &verifyCaptchaErrMsg{ErrMsg: &ErrMsg{Message: "图形验证码错误,请重新输入"}, CaptchaNeeded: true}
	}

	return nil
}
