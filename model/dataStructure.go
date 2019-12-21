package model

import (
	"net/http"
)

// Data 数据是否有效
type Data interface {
	Valid() bool
}

// Request 数据请求类型
type Request struct {
	req   *http.Request
	depth uint32
}

// HTTPReq 返回http请求实例
func (request *Request) HTTPReq() *http.Request {
	return request.req
}

// Depth 返回爬取深度
func (request *Request) Depth() uint32 {
	return request.depth
}

// Valid 请求实例是否有效
func (request *Request) Valid() bool {
	return request.req != nil && request.req.URL != nil
}

// NewRequest 创建一个请求实例
func NewRequest(request *http.Request, deepLen uint32) *Request {
	return &Request{req: request, depth: deepLen}
}

// Response 响应类型
type Response struct {
	rsp   *http.Response
	depth uint32
}

// HTTPRsp 返回http响应实例
func (resp *Response) HTTPRsp() *http.Response {
	return resp.rsp
}

// Depth 返回响应深度
func (resp *Response) Depth() uint32 {
	return resp.depth
}

// Valid 返回响应实例是否有效
func (resp *Response) Valid() bool {
	return resp.rsp != nil && resp.rsp.Body != nil
}

// NewResponse 创建一个响应实例
func NewResponse(response *http.Response, deepLen uint32) *Response {
	return &Response{rsp: response, depth: deepLen}
}

// Item 条目类型
type Item map[string]interface{}

// Valid 返回Item是否有效
func (item Item) Valid() bool {
	return item != nil
}
