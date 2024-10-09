package gpool

import (
	"github.com/pkg/errors"
	"sync"
	"sync/atomic"
)

var workerChanCap = 1

type Config struct {
}

type Pool struct {
	cap         int32
	running     int32
	workers     WorkerQueue
	workerCache sync.Pool // 存储worker
	conf        *Config
	mux         sync.Mutex
}

func NewPool(cap int32, conf *Config) *Pool {
	pool := &Pool{
		cap:     cap,
		workers: newWorkerQueue(10),
		conf:    conf,
	}
	pool.workerCache.New = func() any {
		return &worker{
			pool: pool,
			task: make(chan func(), workerChanCap),
		}
	}
	return pool
}

// Submit 提交任务
// 1. 获取worker
// 2. 将任务交给worker
func (p *Pool) Submit(task func()) error {
	w := p.getWorker()
	if w == nil {
		return errors.New("worker is not exist")
	}
	w.task <- task
	return nil
}

func (p *Pool) getWorker() (w *worker) {
	p.mux.Lock()
	w = p.workers.get()
	if w != nil {
		p.mux.Unlock()
		return
	}
	// 如果没有，则创建
	newWorkerFunc := func() {
		w = p.workerCache.Get().(*worker)
		w.Run()
	}
	if c := p.Cap(); c == -1 || c > p.Running() {
		p.mux.Unlock()
		newWorkerFunc()
	}
	// 当前正在运行的worker大于cap，则直接返回
	p.mux.Unlock()
	return
}

func (p *Pool) Cap() int {
	return int(atomic.LoadInt32(&p.cap))
}

func (p *Pool) Running() int {
	return int(atomic.LoadInt32(&p.running))
}

func (p *Pool) addRunning(k int) {
	atomic.AddInt32(&p.running, int32(k))
}

func (p *Pool) recycleWorker(w *worker) bool {
	p.mux.Lock()
	err := p.workers.put(w)
	if err != nil {
		p.mux.Unlock()
		return false
	}
	p.mux.Unlock()
	return true
}
