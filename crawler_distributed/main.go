package main

import (
	"errors"
	"net/rpc"

	"log"

	"flag"

	"strings"

	"github.com/liuqiqi-Y/qiqiCrawler/crawler/config"
	"github.com/liuqiqi-Y/qiqiCrawler/crawler/engine"
	"github.com/liuqiqi-Y/qiqiCrawler/crawler/scheduler"
	"github.com/liuqiqi-Y/qiqiCrawler/crawler/xcar/parser"
	itemsaver "github.com/liuqiqi-Y/qiqiCrawler/crawler_distributed/persist/client"
	"github.com/liuqiqi-Y/qiqiCrawler/crawler_distributed/rpcsupport"
	worker "github.com/liuqiqi-Y/qiqiCrawler/crawler_distributed/worker/client"
)

var (
	itemSaverHost = flag.String(
		"itemsaver_host", "", "itemsaver host")

	workerHosts = flag.String(
		"worker_hosts", "",
		"worker hosts (comma separated)")
)

func main() {
	flag.Parse()

	itemChan, err := itemsaver.ItemSaver(
		*itemSaverHost)
	if err != nil {
		panic(err)
	}

	pool, err := createClientPool(
		strings.Split(*workerHosts, ","))
	if err != nil {
		panic(err)
	}

	processor := worker.CreateProcessor(pool)

	e := engine.ConcurrentEngine{
		Scheduler:        &scheduler.QueuedScheduler{},
		WorkerCount:      100,
		ItemChan:         itemChan,
		RequestProcessor: processor,
	}

	e.Run(engine.Request{
		Url: "http://www.starter.url.here",
		Parser: engine.NewFuncParser(
			parser.ParseCarList,
			config.ParseCarList),
	})
}

func createClientPool(
	hosts []string) (chan *rpc.Client, error) {
	var clients []*rpc.Client
	for _, h := range hosts {
		client, err := rpcsupport.NewClient(h)
		if err == nil {
			clients = append(clients, client)
			log.Printf("Connected to %s", h)
		} else {
			log.Printf(
				"Error connecting to %s: %v",
				h, err)
		}
	}

	if len(clients) == 0 {
		return nil, errors.New(
			"no connections available")
	}
	out := make(chan *rpc.Client)
	go func() {
		for {
			for _, client := range clients {
				out <- client
			}
		}
	}()
	return out, nil
}
