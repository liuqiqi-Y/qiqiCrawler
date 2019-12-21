package model

import (
	"bytes"
	"fmt"
)

// ErrorType 错误类型
type ErrorType string

// 错误类型常量
const (
	// 下载器错误
	ErrorTypeDownLoader ErrorType = "downloader error"
	// 分祈器错误
	ErrorTypeAnalyzer ErrorType = "analyzer error"
	// 条目处理管道错误
	ErrorTypePipe ErrorType = "pipeline error"
	// 调调度器错误
	ErrorTypeScheduler ErrorType = "scheduler error"
)

// CrawlerError 爬虫错误的接口类型
type CrawlerError interface {
	Type() ErrorType
	Error() string
}

// mycrawLerError 爬虫错误接口的实现
type myCrawlerError struct {
	errType    ErrorType
	errMessage string
	detail     string
}

// Type 返回错误类型
func (err *myCrawlerError) Type() ErrorType {
	return err.errType
}

// Error 返回错误信息
func (err *myCrawlerError) Error() string {
	if err.detail == "" {
		err.getErrDetail()
	}
	return err.detail
}

func (err *myCrawlerError) getErrDetail() {
	var buffer bytes.Buffer
	buffer.WriteString("CrawlerError: ")
	if err.errType != "" {
		buffer.WriteString(string(err.errType))
		buffer.WriteString(":")
	}
	if err.errMessage != "" {
		buffer.WriteString(err.errMessage)
	}
	err.detail = fmt.Sprintf("%s", buffer.String())
	return
}

// NewCrawlerError 返回一个错误实例
func NewCrawlerError(errorType ErrorType, errorMessage string) CrawlerError {
	return &myCrawlerError{
		errType:    errorType,
		errMessage: errorMessage,
	}
}
