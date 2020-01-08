package basemodule

import (
	"errors"
	"fmt"
)

type BaseIF interface {
	ID() MID
	Addr() string
	Score() int
	SetScore(int)
	ScoreCaculator() CalculateScore
	CalledCount() int
	AcceptedCount() int
	CompletedCount() int
	HandlingNumber() int
	Counts() Counts
	Summary() SummartStruct
}

type BaseModuleIF interface {
	BaseIF
	IncrCalledCount()
	IncrAcceptedCount()
	IncrCompletedCount()
	IncrHandlingNumber()
	DecrHandleingNumber()
	Clear()
}

type BaseModule struct {
	mid             MID
	addr            string
	score           int
	scoreCalculator CalculateScore
	calledCount     int
	acceptedCount   int
	completedCount  int
	handlingNumber  int
}

func NewBaseModule(mid MID, scoreCalculator CalculateScore) (BaseModuleIF, error) {
	parts, err := SplitMID(mid)
	if err != nil {

		return nil, errors.New(fmt.Sprintf("illegal MID: %s, %s", mid, err.Error()))
	}
}
