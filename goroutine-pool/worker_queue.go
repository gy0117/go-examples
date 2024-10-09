package gpool

type WorkerQueue interface {
	get() *worker
	put(w *worker) error
}

type workerQueue struct {
	items []*worker
}

var _ WorkerQueue = (*workerQueue)(nil)

func newWorkerQueue(size int) *workerQueue {
	return &workerQueue{
		items: make([]*worker, 0, size),
	}
}

func (wq *workerQueue) get() *worker {
	n := len(wq.items)
	if n == 0 {
		return nil
	}
	w := wq.items[0]
	wq.items = wq.items[1:]
	return w
}

func (wq *workerQueue) put(w *worker) error {
	wq.items = append(wq.items, w)
	return nil
}
