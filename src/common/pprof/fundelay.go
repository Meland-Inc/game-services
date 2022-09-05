package profile

import (
	"time"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
)

type FunDelay struct {
	Name string
	T0   int64
	T9   int64
}

type FunDelays struct {
	Chans       chan *FunDelay
	CheckChan   chan int32
	ResetChan   chan struct{}
	StopChan    chan chan struct{}
	Maps        map[string]int64
	Nums        map[string]int64
	Detailes    map[string]map[int32]int64
	DetaileNums map[string]int32

	bStop bool
}

var funs *FunDelays

func initDataContainers() {
	if funs != nil {
		return
	}
	funs = new(FunDelays)
	funs.Chans = make(chan *FunDelay, 51200)
	funs.CheckChan = make(chan int32, 1)
	funs.ResetChan = make(chan struct{}, 1)
	funs.Maps = make(map[string]int64)
	funs.Nums = make(map[string]int64)
	funs.Detailes = make(map[string]map[int32]int64)
	funs.DetaileNums = make(map[string]int32)

	funs.StopChan = make(chan chan struct{})

}

func StartSpans() {
	initDataContainers()
	go func() {
		defer func() {
			err := recover()
			if err != nil {
				go StartSpans()
			}
		}()

		for {
			select {
			case f := <-funs.Chans:
				funs.Maps[f.Name] += (f.T9 - f.T0)
				funs.Nums[f.Name]++
				if funs.Detailes[f.Name] == nil {
					funs.Detailes[f.Name] = make(map[int32]int64)
					funs.DetaileNums[f.Name] = 1
				}
				funs.Detailes[f.Name][funs.DetaileNums[f.Name]] = (f.T9 - f.T0)
				funs.DetaileNums[f.Name]++
			case <-funs.CheckChan:
				showSpans()
			case <-funs.ResetChan:
				resetSpans()
			case stopFinished := <-funs.StopChan:
				close(funs.StopChan)
				close(funs.Chans)
				close(funs.CheckChan)
				close(funs.ResetChan)
				stopFinished <- struct{}{}
				return
			}
		}
	}()
}

func StartSpan(name string) *FunDelay {
	return &FunDelay{
		Name: name,
		T0:   time.Now().UnixMicro(),
	}
}

func (f *FunDelay) Finish() {
	if funs.bStop {
		return
	}
	f.T9 = time.Now().UnixMicro()
	if len(funs.Chans) < 50000 {
		funs.Chans <- f
	} else {
		go func() {
			funs.Chans <- f
		}()
	}
}

func ShowSpans() {
	funs.CheckChan <- 1
}

func showSpans() {
	serviceLog.Info("-----------------------------------------")
	for name, allT := range funs.Maps {
		count := funs.Nums[name]
		serviceLog.Info("Func(%s) Need(%f)us, callNum(%d)", name, float64(allT)/float64(count), count)
	}
	serviceLog.Info("-----------------------------------------")

	resetSpans()

}

func resetSpans() {
	funs.Maps = make(map[string]int64)
	funs.Nums = make(map[string]int64)
	funs.Detailes = make(map[string]map[int32]int64)
	funs.DetaileNums = make(map[string]int32)
}

func Stop() {
	funs.bStop = true
	stopDone := make(chan struct{}, 1)
	funs.StopChan <- stopDone
	<-stopDone
	funs = nil
	serviceLog.Info("Spans Stop")
}
