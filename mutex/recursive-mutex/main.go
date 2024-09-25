package main

import (
	"fmt"
	"github.com/kortschak/goroutine"
	"sync"
	"sync/atomic"
	"time"
)

// Mutex扩展：可重入锁

type RecursiveMutex struct {
	sync.Mutex
	// gid，哪个g持有了当前锁
	owner int64
	// 重入次数
	recursiveCount int64
}

func (m *RecursiveMutex) Lock() {
	// 判断gid是否与当前G的id一致，如果一致，递归次数++，return；如果不一致，则加锁
	gid := getGid()
	if atomic.LoadInt64(&m.owner) == gid {
		atomic.AddInt64(&m.recursiveCount, 1)
		return
	}

	m.Lock()
	atomic.StoreInt64(&m.owner, gid)
	atomic.StoreInt64(&m.recursiveCount, 1)
}

func (m *RecursiveMutex) UnLock() {
	gid := getGid()
	if atomic.LoadInt64(&m.owner) != gid {
		panic("只允许加锁的G是否锁")
	}
	if atomic.AddInt64(&m.recursiveCount, -1) != 0 {
		return
	}
	atomic.StoreInt64(&m.owner, -1)
	m.Unlock()
}

func getGid() int64 {
	return goroutine.ID()
}

func main() {
	gid := getGid()
	fmt.Println("gid: ", gid)

	go func() {
		gid := getGid()
		fmt.Println("gid func1: ", gid)
	}()

	go func() {
		gid := getGid()
		fmt.Println("gid func2: ", gid)
	}()

	time.Sleep(time.Second)
}
