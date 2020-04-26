package handler

import (
	"bytes"
	"net/http"
	"net/url"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/skip2/go-qrcode"
)

// text=http%3A%2F%2Fstatic.chat186.com%2Fdist%2Fstandalone.html%3Feid%3Dbjstn2pfua60sjh9jae0
func (s *IMService) QrCode(ctx echo.Context) error {
	text := ctx.QueryParam("text")
	if text == "" {
		return invalidParameterResp(ctx, "empty text")
	}

	link, err := url.QueryUnescape(text)
	if err != nil {
		return invalidParameterResp(ctx, err.Error())
	}

	var png []byte
	png, err = qrcode.Encode(link, qrcode.Medium, 256)
	if err != nil {
		return invalidParameterResp(ctx, err.Error())
	}

	w := ctx.Response()
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	w.Header().Set("Content-Type", "image/png")

	r := ctx.Request()
	http.ServeContent(w, r, "chat-link.png", time.Time{}, bytes.NewReader(png))
	return nil
}
