package main

import (
	"log"
	"sync"
	"time"
)

func main() {
	Wg()
}

func Wg() {
	taskNum := 10
	var wg sync.WaitGroup
	dataCh := make(chan any)
	stopCh := make(chan struct{}, 1)
	resp := make([]any, 0, taskNum)

	// 读
	go func() {
		for data := range dataCh {
			resp = append(resp, data)
		}
		stopCh <- struct{}{}
	}()

	log.Println("start")
	for i := 0; i < taskNum; i++ {
		wg.Add(1)
		go func(ch chan any) {
			defer wg.Done()
			ch <- time.Now().UnixNano()
		}(dataCh)
	}

	wg.Wait()
	close(dataCh)
	log.Println("done.")

	select {
	case <-stopCh:
		log.Println(resp)
	default:

	}

	//log.Println(resp)
}

// 主g充当读g，新开一个写g
func Wg2() {
	taskNum := 10
	dataCh := make(chan any)

	go func() {
		var wg sync.WaitGroup
		for i := 0; i < taskNum; i++ {
			wg.Add(1)
			go func(ch chan any) {
				defer wg.Done()
				ch <- time.Now().UnixNano()
			}(dataCh)
		}
		wg.Wait()
		close(dataCh)
	}()

	resp := make([]any, 0, taskNum)
	for data := range dataCh {
		resp = append(resp, data)
	}
	log.Println(resp)
}
