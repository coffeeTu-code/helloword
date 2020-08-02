package pipeline

import (
	"errors"
	"sync/atomic"
)

type Filter interface {
	ProcessWithMetrics(ctx interface{}) (modelStatus string, time int64, err error)
	Process(ctx interface{}) (modelStatus string, err error)
	PipelineName() (name string)
}

//---------------------------------------

type AtomicInt int64

func (a *AtomicInt) Add(i int64) {
	atomic.AddInt64((*int64)(a), i)
}

func (a *AtomicInt) Val() int64 {
	return *(*int64)(a)
}

//--------------------------------------

type WallTimePipeline struct {
	Name        string
	Filters     []Filter
	TimeElapsed []AtomicInt
}

//ProcessWithTime 执行pipeline，并记录各模块执行时间
func (wtp *WallTimePipeline) ProcessWithTime(ctx interface{}) error {
	if ctx == nil {
		return errors.New("WallTimePipeline: function ProcessWithTime / ctx is nil")
	}
	if wtp == nil || len(wtp.Filters) == 0 {
		return errors.New("WallTimePipeline: function ProcessWithTime / Filters is nil")
	}

	wtp.TimeElapsed = make([]AtomicInt, len(wtp.Filters))
	for i, filter := range wtp.Filters {
		_, processTimeElapsed, err := filter.ProcessWithMetrics(ctx)

		wtp.TimeElapsed[i].Add(processTimeElapsed)
		if err != nil {
			return errors.New("WallTimePipeline: " + err.Error())
		}
	}
	return nil
}

func (wtp *WallTimePipeline) TotalTime() (total int64) {
	if wtp == nil {
		return 0
	}

	for i, _ := range wtp.TimeElapsed {
		total += wtp.TimeElapsed[i].Val()
	}
	return total
}
