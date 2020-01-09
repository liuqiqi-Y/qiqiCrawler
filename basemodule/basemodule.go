package basemodule

import (
	"errors"
	"fmt"
	"sync/atomic"

	"github.com/liuqiqi-Y/qiqiCrawler/util"
)

type BaseIF interface {
	ID() MID
	Addr() string
	Score() uint64
	SetScore(uint64)
	ScoreCalculator() CalculateScore
	CalledCount() uint64
	AcceptedCount() uint64
	CompletedCount() uint64
	HandlingNumber() uint64
	Counts() Counts
	Summary() SummaryStruct
}

type BaseModuleIF interface {
	BaseIF
	IncrCalledCount()
	IncrAcceptedCount()
	IncrCompletedCount()
	IncrHandlingNumber()
	DecrHandlingNumber()
	Clear()
}

type BaseModule struct {
	mid             MID
	addr            string
	score           uint64
	scoreCalculator CalculateScore
	calledCount     uint64
	acceptedCount   uint64
	completedCount  uint64
	handlingNumber  uint64
}

func NewBaseModule(id MID, Calculator CalculateScore) (BaseModuleIF, error) {
	parts, err := SplitMID(id)
	if err != nil {
		util.Trace.Printf("MID is %s", id)
		return nil, errors.New(fmt.Sprintf("illegal MID: %s, %s", id, err.Error()))
	}
	return &BaseModule{
		mid:             id,
		addr:            parts[2],
		scoreCalculator: Calculator,
	}, nil
}
func (m *BaseModule) ID() MID {
	return m.mid
}

func (m *BaseModule) Addr() string {
	return m.addr
}

func (m *BaseModule) Score() uint64 {
	return atomic.LoadUint64(&m.score)
}

func (m *BaseModule) SetScore(score uint64) {
	atomic.StoreUint64(&m.score, score)
}

func (m *BaseModule) ScoreCalculator() CalculateScore {
	return m.scoreCalculator
}

func (m *BaseModule) CalledCount() uint64 {
	return atomic.LoadUint64(&m.calledCount)
}

func (m *BaseModule) AcceptedCount() uint64 {
	return atomic.LoadUint64(&m.acceptedCount)
}

func (m *BaseModule) CompletedCount() uint64 {
	return atomic.LoadUint64(&m.completedCount)
}

func (m *BaseModule) HandlingNumber() uint64 {
	return atomic.LoadUint64(&m.handlingNumber)
}

func (m *BaseModule) Counts() Counts {
	return Counts{
		CalledCount:    atomic.LoadUint64(&m.calledCount),
		AcceptedCount:  atomic.LoadUint64(&m.acceptedCount),
		CompletedCount: atomic.LoadUint64(&m.completedCount),
		HandlingNumber: atomic.LoadUint64(&m.handlingNumber),
	}
}

func (m *BaseModule) Summary() SummaryStruct {
	counts := m.Counts()
	return SummaryStruct{
		ID:        m.ID(),
		Called:    counts.CalledCount,
		Accepted:  counts.AcceptedCount,
		Completed: counts.CompletedCount,
		Handling:  counts.HandlingNumber,
		Extra:     nil,
	}
}

func (m *BaseModule) IncrCalledCount() {
	atomic.AddUint64(&m.calledCount, 1)
}

func (m *BaseModule) IncrAcceptedCount() {
	atomic.AddUint64(&m.acceptedCount, 1)
}

func (m *BaseModule) IncrCompletedCount() {
	atomic.AddUint64(&m.completedCount, 1)
}

func (m *BaseModule) IncrHandlingNumber() {
	atomic.AddUint64(&m.handlingNumber, 1)
}

func (m *BaseModule) DecrHandlingNumber() {
	atomic.AddUint64(&m.handlingNumber, ^uint64(0))
}

func (m *BaseModule) Clear() {
	atomic.StoreUint64(&m.calledCount, 0)
	atomic.StoreUint64(&m.acceptedCount, 0)
	atomic.StoreUint64(&m.completedCount, 0)
	atomic.StoreUint64(&m.handlingNumber, 0)
}
