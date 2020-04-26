package handler

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"path/filepath"
	"time"

	"github.com/labstack/echo/v4"

	"bitbucket.org/forfd/custm-chat/webim/common"
	"bitbucket.org/forfd/custm-chat/webim/common/logging"
	"bitbucket.org/forfd/custm-chat/webim/models"
)

var (
	defaultExpires     = 24 * time.Hour
	defaultMaxFileSize = int64(5 * math.Pow10(7)) // 50M
)

type UploadResp struct {
	ExpiresIn int    `json:"expires_in"`
	FileURL   string `json:"file_url"`
}

type AdminUploadResp struct {
	Success  bool   `json:"success"`
	PhotoKey string `json:"photo_key"`
	PhotoURL string `json:"photo_url"`
}

// Upload ...
// POST /api/v1/upload?ent_id=xxxxx&type=file/photo
func (s *IMService) Upload(ctx echo.Context) (err error) {
	entID := ctx.QueryParam("ent_id")
	fileType := ctx.QueryParam("type")
	if entID == "" || fileType == "" {
		return invalidParameterResp(ctx, "ent_id/type invalid")
	}

	if fileType != "photo" && fileType != "file" {
		return invalidParameterResp(ctx, "unsupported file type")
	}

	var expireAt = time.Time{}
	if fileType == "file" {
		expireAt = time.Now().Add(defaultExpires)
	}

	extra, msg := s.upload(ctx, entID, expireAt)
	if msg != nil {
		return invalidParameterResp(ctx, msg.Message)
	}

	resp := &UploadResp{ExpiresIn: -1, FileURL: extra.Name}
	if fileType == "file" {
		resp.ExpiresIn = int(defaultExpires.Minutes())
	}

	return jsonResponse(ctx, resp)
}

// POST /upload_img/
func (s *IMService) UploadImg(ctx echo.Context) (err error) {
	entID := ctx.FormValue("ent_id")
	base64 := ctx.FormValue("content")
	name := ctx.FormValue("name")
	if entID == "" || base64 == "" {
		return invalidParameterResp(ctx, "ent_id/base64 invalid")
	}

	if name == "" {
		name = fmt.Sprintf("%d", time.Now().Unix())
	}

	var expireAt = time.Time{}
	extra, msg := s.uploadImgBytes(name, entID, base64, expireAt)
	if msg != nil {
		return invalidParameterResp(ctx, msg.Message)
	}

	resp := &UploadResp{ExpiresIn: -1, FileURL: extra.Name}
	return jsonResponse(ctx, resp)
}

func (s *IMService) AdminUpload(ctx echo.Context) (err error) {
	url, msg := s.uploadV1(ctx)
	if msg != nil {
		return jsonResponse(ctx, msg)
	}

	_, file := filepath.Split(url)
	return jsonResponse(ctx, &AdminUploadResp{Success: true, PhotoKey: file, PhotoURL: url})
}

func (s *IMService) upload(ctx echo.Context, entID string, expireAt time.Time) (res *models.FileExtra, msg *ErrMsg) {
	extra := &models.FileExtra{
		EntID: entID,
	}

	file, err := ctx.FormFile("file")
	if err != nil {
		return nil, &ErrMsg{Code: common.UploadFileErr, Message: err.Error()}
	}

	if file.Size > defaultMaxFileSize {
		return nil, &ErrMsg{Code: common.UploadFileErr, Message: "file too large"}
	}

	src, err := file.Open()
	if err != nil {
		return nil, &ErrMsg{Code: common.UploadFileErr, Message: err.Error()}
	}
	defer func() {
		if closeErr := src.Close(); closeErr != nil {
			log.Logger.Warnf("close file error: %v", closeErr)
		}
	}()

	var fileType string
	if v, ok := file.Header[echo.HeaderContentType]; ok {
		if len(v) > 0 {
			fileType = v[0]
		}
	}

	bs, err := CompressImg(fileType, src)
	if err != nil {
		log.Logger.Warn(" Compress Img error: ", err)
	}

	var fileReader io.Reader = src
	if len(bs) > 0 {
		fileReader = bytes.NewBuffer(bs)
	}

	fileName := fmt.Sprintf("%s/%s/%s", entID, "files", file.Filename)
	url, err := s.uploader.Upload(fileName, fileReader, expireAt)
	if err != nil {
		log.Logger.Warn("upload file(size: ", file.Size, ", name: ", file.Filename, ") error: ", err)
		return nil, &ErrMsg{Code: common.UploadFileErr, Message: err.Error()}
	}

	now := time.Now().UTC()
	extra.Name = url
	extra.Size = int(file.Size)
	extra.UploadTime = now
	extra.ExpireAt = expireAt
	extra.Type = fileType
	if expireAt.IsZero() {
		extra.ExpireAt = now.Add(24 * time.Hour * 365 * 100)
	}

	return extra, nil
}

func (s *IMService) uploadV1(ctx echo.Context) (url string, errMsg *ErrMsg) {
	file, err := ctx.FormFile("file")
	if err != nil {
		return "", &ErrMsg{Code: common.UploadFileErr, Message: err.Error()}
	}

	if file.Size > defaultMaxFileSize {
		return "", &ErrMsg{Code: common.UploadFileErr, Message: "file too large"}
	}

	src, err := file.Open()
	if err != nil {
		return "", &ErrMsg{Code: common.UploadFileErr, Message: err.Error()}
	}
	defer func() {
		if closeErr := src.Close(); closeErr != nil {
			log.Logger.Warnf("close file error: %v", closeErr)
		}
	}()

	fileName := fmt.Sprintf("%s/%s-%d-%s", "files", "file", time.Now().Unix(), file.Filename)
	url, err = s.uploader.Upload(fileName, src, time.Time{})
	if err != nil {
		log.Logger.Warn("upload file(size: ", file.Size, ", name: ", file.Filename, ") error: ", err)
		return "", &ErrMsg{Code: common.UploadFileErr, Message: err.Error()}
	}

	return
}

func (s *IMService) uploadImgBytes(fileName, entID, base64 string, expireAt time.Time) (res *models.FileExtra, msg *ErrMsg) {
	extra := &models.FileExtra{
		EntID: entID,
	}

	imgBytes, err := base64ToImgBytes(base64)
	if err != nil {
		return nil, &ErrMsg{Code: common.UploadFileErr, Message: err.Error()}
	}
	fileReader := bytes.NewBuffer(imgBytes)

	fileName = fmt.Sprintf("%s/%s/%s", entID, "files", fileName)
	url, err := s.uploader.Upload(fileName, fileReader, expireAt)
	if err != nil {
		return nil, &ErrMsg{Code: common.UploadFileErr, Message: err.Error()}
	}

	now := time.Now().UTC()
	extra.Name = url
	extra.UploadTime = now
	extra.ExpireAt = expireAt
	if expireAt.IsZero() {
		extra.ExpireAt = now.Add(24 * time.Hour * 365 * 100)
	}

	return extra, nil
}
