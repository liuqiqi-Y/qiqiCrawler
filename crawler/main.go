package main

import (
	"github.com/liuqiqi-Y/qiqiCrawler/config"
	"github.com/liuqiqi-Y/qiqiCrawler/engine"
	"github.com/liuqiqi-Y/qiqiCrawler/persist"
	"github.com/liuqiqi-Y/qiqiCrawler/scheduler"
	"github.com/liuqiqi-Y/qiqiCrawler/xcar/parser"
)

func main() {
	itemChan, err := persist.ItemSaver(
		config.ElasticIndex)
	if err != nil {
		panic(err)
	}

	e := engine.ConcurrentEngine{
		Scheduler:        &scheduler.QueuedScheduler{},
		WorkerCount:      100,
		ItemChan:         itemChan,
		RequestProcessor: engine.Worker,
	}

	e.Run(engine.Request{
		Url: "http://www.starter.url.here",
		Parser: engine.NewFuncParser(
			parser.ParseCarList,
			config.ParseCarList),
	})
}
