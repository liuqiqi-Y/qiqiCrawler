package basemodule

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/liuqiqi-Y/qiqiCrawler/util"
)

type MID string

var midTemplate = "%s%d|%s"

type Counts struct {
	// CalledCount 代表调用计数。
	CalledCount uint64
	// AcceptedCount 代表接受计数。
	AcceptedCount uint64
	// CompletedCount 代表成功完成计数。
	CompletedCount uint64
	// HandlingNumber 代表实时处理数。
	HandlingNumber uint64
}

// SummaryStruct 代表组件摘要结构的类型。
type SummaryStruct struct {
	ID        MID         `json:"id"`
	Called    uint64      `json:"called"`
	Accepted  uint64      `json:"accepted"`
	Completed uint64      `json:"completed"`
	Handling  uint64      `json:"handling"`
	Extra     interface{} `json:"extra,omitempty"`
}

const (
	// TYPE_DOWNLOADER 代表下载器。
	TYPE_DOWNLOADER = "downloader"
	// TYPE_ANALYZER 代表分析器。
	TYPE_ANALYZER = "analyzer"
	// TYPE_PIPELINE 代表条目处理管道。
	TYPE_PIPELINE = "pipeline"
)

var legalTypeLetterMap = map[string]string{
	TYPE_DOWNLOADER: "D",
	TYPE_ANALYZER:   "A",
	TYPE_PIPELINE:   "P",
}

var legalLetterTypeMap = map[string]string{
	"D": TYPE_DOWNLOADER,
	"A": TYPE_ANALYZER,
	"P": TYPE_PIPELINE,
}

func GenMID(moduletype string, sn int, addr net.Addr) (MID, error) {
	letter, ok := legalTypeLetterMap[moduletype]
	if !ok {
		util.Trace.Printf("Module type is: %s", moduletype)
		return "", errors.New(fmt.Sprintf("illeagal module type: %s", moduletype))
	}
	var midStr string
	if addr == nil {
		midStr = fmt.Sprintf(midTemplate, letter, sn, "")
		midStr = midStr[:len(midStr)-1]
	} else {
		midStr = fmt.Sprintf(midTemplate, letter, sn, addr.String())
	}
	return MID(midStr), nil
}

func SplitMID(mid MID) ([]string, error) {
	if mid == "" {
		util.Trace.Println("Module ID is NULL")
		return nil, errors.New("Module Id can not be NULL")
	}
	var snStr, addr string
	midStr := string(mid)
	letter := midStr[:1]
	if _, exist := legalLetterTypeMap[letter]; !exist {
		util.Trace.Printf("Module type letter is %s", letter)
		return nil, errors.New(fmt.Sprintf("illegal module type letter: %s", letter))
	}
	snAndAddr := midStr[1:]
	index := strings.LastIndex(snAndAddr, "|")
	if index < 0 {
		snStr = snAndAddr
		if !legalSN(snStr) {
			util.Trace.Printf("module SN is %s", snStr)
			return nil, errors.New(fmt.Sprintf("illegal module SN: %s", snStr))
		}
	} else {
		snStr = snAndAddr[:index]
		if !legalSN(snStr) {
			util.Trace.Printf("module SN is %s", snStr)
			return nil, errors.New(fmt.Sprintf("illegal module SN: %s", snStr))
		}
		addr = snAndAddr[index+1:]
		index = strings.LastIndex(addr, ":")
		if index <= 0 {
			util.Trace.Printf("module address is %s", snStr)
			return nil, errors.New(fmt.Sprintf("illegal module address: %s", addr))
		}
		ipStr := addr[:index]
		if ip := net.ParseIP(ipStr); ip == nil {
			util.Trace.Printf("module IP is %s", snStr)
			return nil, errors.New(fmt.Sprintf("illegal module IP: %s", addr))
		}
		portStr := addr[index+1:]
		if _, err := strconv.ParseUint(portStr, 10, 0); err != nil {
			util.Trace.Printf("module port is %s", snStr)
			return nil, errors.New(fmt.Sprintf("illegal module port: %s", addr))
		}
	}
	return []string{letter, snStr, addr}, nil
}
func legalSN(snStr string) bool {
	_, err := strconv.ParseUint(snStr, 10, 0)
	if err != nil {
		return false
	}
	return true
}
func LegalMID(mid MID) bool {
	if _, err := SplitMID(mid); err == nil {
		return true
	}
	return false
}
