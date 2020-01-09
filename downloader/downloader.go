package downloader

import (
	"fmt"
	"net/http"

	"github.com/liuqiqi-Y/qiqiCrawler/basemodule"
	"github.com/liuqiqi-Y/qiqiCrawler/model"
)

type DownloaderIF interface {
	basemodule.BaseIF
	// Download 会根据请求获取内容并返回响应。
	Download(req *model.Request) (*model.Response, error)
}
type Downloader struct {
	// stub.ModuleInternal 代表组件基础实例。
	basemodule.BaseModuleIF
	// httpClient 代表下载用的HTTP客户端。
	httpClient http.Client
}

func New(mid basemodule.MID, client *http.Client, scoreCalculator basemodule.CalculateScore) (DownloaderIF, error) {
	moduleBase, err := basemodule.NewBaseModule(mid, scoreCalculator)
	if err != nil {
		return nil, err
	}
	if client == nil {
		return nil, genParameterError("nil http client")
	}
	return &Downloader{
		BaseModuleIF: moduleBase,
		httpClient:   *client,
	}, nil
}

func (downloader *Downloader) DownLoad(req *model.Request) (*model.Response, error) {
	downloader.BaseModuleIF.IncrHandlingNumber()
	defer downloader.BaseModuleIF.DecrHandlingNumber()
	downloader.BaseModuleIF.IncrCalledCount()
	if req == nil {
		return nil, fmt.Errorf("")
	}
}
