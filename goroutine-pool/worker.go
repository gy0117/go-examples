package gpool

type worker struct {
	// 这个worker所属的池子
	pool *Pool
	task chan func()
}

func (w *worker) Run() {
	w.pool.addRunning(1)
	// worker运行
	go func() {
		defer func() {
			w.pool.addRunning(-1)
			// 退出for循环后，此worker就不再运行了，就可以放回对象池了
			w.pool.workerCache.Put(w)
			if p := recover(); p != nil {
				// handle panic
			}
		}()

		for f := range w.task {
			if f == nil {
				return
			}
			f()
			if ok := w.pool.recycleWorker(w); !ok {
				return
			}
		}
	}()
}
