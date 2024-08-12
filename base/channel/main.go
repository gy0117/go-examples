package main

import (
	"log"
	"sync"
)

func testCase() {
	ch := make(chan int, 10)

	for i := 0; i < 10; i++ {
		ch <- i
	}
	// 上述for循环执行10次之后，就不存在ch <- i了，只剩下 <-ch，会报错
	// 关闭后，ok为false，退出下面的循环
	close(ch)
	for {
		v, ok := <-ch
		if !ok {
			break
		}
		log.Println(v)
	}
}

func main() {
	counter := 0
	counterChan := make(chan int)

	// 创建等待组，用于等待所有goroutine完成
	var wg sync.WaitGroup

	// 10个goroutine并发递增计数器
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			// 递增计数器
			counterChan <- 1
			wg.Done()
		}()
	}

	// 等待所有goroutine完成
	go func() {
		wg.Wait()
		close(counterChan) // 关闭通道，告诉接收方，我这里数据已经发送结束了
	}()

	// 从通道中读取计数器增量，并累加到counter
	for increment := range counterChan {
		counter += increment
	}
	log.Println("Final counter value:", counter)
}
