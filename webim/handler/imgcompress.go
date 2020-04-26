package handler

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image/jpeg"
	"image/png"
	"io"
)

const defaultQuality = 30

func CompressImg(contentType string, reader io.Reader) (bs []byte, err error) {
	switch contentType {
	case "image/jpeg":
		return jpegCompress(reader)
	case "image/png":
		return pngCompress(reader)
	default:
		return nil, fmt.Errorf("content-type not support")
	}
}

func jpegCompress(reader io.Reader) (bs []byte, err error) {
	img, err := jpeg.Decode(reader)
	if err != nil {
		return nil, err
	}

	var buf = new(bytes.Buffer)
	err = jpeg.Encode(buf, img, &jpeg.Options{Quality: defaultQuality})
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func pngCompress(reader io.Reader) (bs []byte, err error) {
	img, err := png.Decode(reader)
	if err != nil {
		return nil, err
	}

	encoder := png.Encoder{CompressionLevel: png.BestCompression}

	var buf = new(bytes.Buffer)
	err = encoder.Encode(buf, img)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func base64ToImgBytes(base64Str string) ([]byte, error) {
	unbased, err := base64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		return nil, err
	}
	return unbased, nil
}
