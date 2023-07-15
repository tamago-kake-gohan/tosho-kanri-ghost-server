package utils

import (
	"net/http"
)

func GetRequestBody(r *http.Request) []byte {
	len := r.ContentLength
	rawBody := make([]byte, len) // Content-Length と同じサイズの byte 配列を用意
	r.Body.Read(rawBody)         // byte 配列にリクエストボディを読み込む
	return rawBody
}
