package main

import (
	"context"
	"sync"
	"time"
)

// 考题
// 要求实现一个map
// 1. 面向高并发
// 2. 只存在插入和查询操作 O(1)
// 3. 查询时，若key存在，直接返回val；若key不存在，阻塞直到key val 对被放入后，获取val 返回；等待指定时长仍未放入，返回超时错误
// 4. 写出真实代码，不能有死锁或者panic风险

type MyChan struct {
	sync.Once
	ch chan struct{}
}

func NewMyChan() *MyChan {
	return &MyChan{
		ch: make(chan struct{}),
	}
}

func (ch *MyChan) Close() {
	ch.Do(func() {
		close(ch.ch)
	})
}

type MyConcurrentMap struct {
	sync.Mutex
	mp    map[int]int
	keyCh map[int]*MyChan
}

func NewMyConcurrentMap() *MyConcurrentMap {
	return &MyConcurrentMap{
		mp:    make(map[int]int),
		keyCh: make(map[int]*MyChan),
	}
}

func (m *MyConcurrentMap) Put(k, v int) {
	m.Lock()
	defer m.Unlock()
	m.mp[k] = v

	ch, ok := m.keyCh[k]
	// 不存在，就没人读，直接return
	if !ok {
		return
	}
	// 说明有人正在阻塞住，等待读
	// 可能存在多个人正等待读，因此不能 ch <- struct{}{}
	// 只是一个close方法，可能会执行多次，会panic
	//close(ch)
	ch.Close()
}

func (m *MyConcurrentMap) Get(k int, maxWaitingDuration time.Duration) (int, error) {
	m.Lock()
	v, ok := m.mp[k]
	if ok {
		m.Unlock()
		return v, nil
	}

	// 不存在
	// 1. 阻塞住，等待插入的时候，再次唤醒
	ch, ok := m.keyCh[k]
	if !ok {
		ch = NewMyChan()
		m.keyCh[k] = ch
	}
	ctx, cancel := context.WithTimeout(context.Background(), maxWaitingDuration)
	defer cancel()
	m.Unlock()

	// 阻塞住
	select {
	case <-ctx.Done():
		return -1, ctx.Err()
	case <-ch.ch:
	}

	// 此时也可能有人正在写，因此需要加锁
	m.Lock()
	v = m.mp[k]
	m.Unlock()
	return v, nil
}
