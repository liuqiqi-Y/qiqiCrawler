package model

import (
	"net/http"
)

// Request 数据请求类型
type Request struct {
	httpReq *http.Request
}
