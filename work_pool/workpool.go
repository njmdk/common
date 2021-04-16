package workpool

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"go.uber.org/zap"

	"github.com/njmdk/common/logger"
	"github.com/njmdk/common/utils"
)

type WorkPool struct {
	cap       int32
	current   int32
	tasks     chan func()
	postTasks chan func()
	close     chan struct{}
	wg        *sync.WaitGroup

	chanPool chan *Worker

	log *logger.Logger
}

// 传入参数暂时无用
func NewWorkPool(maxPool int32, log *logger.Logger) *WorkPool {
	if maxPool<=0{
		maxPool = 10000
	}
	wp := &WorkPool{
		tasks:    make(chan func(), maxPool),
		cap:      maxPool,
		wg:       &sync.WaitGroup{},
		chanPool: make(chan *Worker, maxPool),
		close:    make(chan struct{}),
		log:      log,
	}
	wp.postTasks = wp.tasks
	return wp
}

func (this_ *WorkPool) Post(f func()) {
	select {
	case <-this_.close:
	case this_.postTasks <- f:
	}
}

func (this_ *WorkPool) Run(panicFunc func(i interface{})) {
	if panicFunc == nil {
		panicFunc = func(e interface{}) {
			if this_.log != nil {
				this_.log.Error("workpool panic", zap.Any("panic info", e))
			} else {
				fmt.Println("workpool panic", e)
			}
		}
	}
	this_.chanPool <- this_.createWorkerOnce()
	this_.wg.Add(1)
	utils.SafeGO(panicFunc, func() {
		defer this_.wg.Done()
		for v := range this_.tasks {
			this_.dealOneTask(v)
		}
		close(this_.close)
		for oneChan := range this_.chanPool {
			if oneChan.close {
				select {
				case v, ok := <-oneChan.f:
					if ok {
						if v != nil {
							this_.createWorker().f <- v
						}
						close(oneChan.f)
					}
				default:
					close(oneChan.f)
				}
			}
		}
	})
}

func (this_ *WorkPool) dealOneTask(f func()) {
	for {
		select {
		case oneChan := <-this_.chanPool:
			if oneChan.close {
				select {
				case v, ok := <-oneChan.f:
					if ok {
						if v != nil {
							this_.createWorker().f <- v
						}
						close(oneChan.f)
					}
				default:
					close(oneChan.f)
				}
				continue
			}
			oneChan.f <- f
		default:
			this_.createWorker().f <- f
		}
		return
	}
}

type Worker struct {
	f     chan func()
	close bool
}

func (this_ *WorkPool) createWorkerOnce() *Worker {
	worker := &Worker{
		f: make(chan func(), 1),
	}
	this_.wg.Add(1)
	panicFunc := func(i interface{}) {
		this_.log.Error("createWorker panic", zap.Any("panic info", i))
	}
	atomic.AddInt32(&this_.current, 1)
	utils.SafeGO(panicFunc, func() {
		defer func() {
			closeNum := atomic.AddInt32(&this_.current, -1)
			if closeNum == 0 {
				close(this_.chanPool)
			} else {
				this_.log.Debug("close work chan", zap.Int32("closeNum", closeNum))
			}
			this_.wg.Done()
		}()
		for {
			select {
			case v, ok := <-worker.f:
				if !ok {
					return
				}
				func() {
					defer utils.Recover(panicFunc)
					v()
				}()
				this_.chanPool <- worker
			case <-this_.close:
				closeWorker(worker)
				this_.chanPool <- worker
				return
			}
		}
	})
	return worker
}

func (this_ *WorkPool) createWorker() *Worker {
	worker := &Worker{
		f: make(chan func(), 1),
	}
	this_.wg.Add(1)
	atomic.AddInt32(&this_.current, 1)
	panicFunc := func(i interface{}) {
		this_.log.Error("createWorker panic", zap.Any("panic info", i))
	}
	utils.SafeGO(panicFunc, func() {
		defer func() {
			closeNum := atomic.AddInt32(&this_.current, -1)
			if closeNum == 0 {
				close(this_.chanPool)
			} else {
				this_.log.Debug("close work chan", zap.Int32("closeNum", closeNum))
			}
			this_.wg.Done()
		}()

		timerDuration := time.Second * 30
		closeTimer := time.NewTimer(timerDuration)
		for {
			select {
			case v, ok := <-worker.f:
				if !ok {
					return
				}
				func() {
					defer utils.Recover(panicFunc)
					v()
				}()

				this_.chanPool <- worker
				closeTimer.Reset(time.Second * 30)
			case <-this_.close:
				closeWorker(worker)
				this_.chanPool <- worker
				return
			case <-closeTimer.C:
				closeWorker(worker)
				this_.chanPool <- worker
				return
			}
		}
	})
	return worker
}

func closeWorker(w *Worker) {
	select {
	case v := <-w.f:
		if v != nil {
			v()
		}
	default:
	}
	w.close = true
}

func (this_ *WorkPool) Stop() {
	this_.Close()
	this_.WaitClosed()
}

func (this_ *WorkPool) Close() {
	this_.postTasks = nil
	close(this_.tasks)
}

func (this_ *WorkPool) WaitClosed() {
	this_.wg.Wait()
}
