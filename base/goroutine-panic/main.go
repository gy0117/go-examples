package main

import (
	"fmt"
	"time"
)

func main() {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Recovered from panic:", r)
			}
		}()
		panic("something went wrong")
	}()

	// 主 goroutine 继续执行，不受子 goroutine 影响
	time.Sleep(time.Second)
	fmt.Println("Main goroutine completed")

	time.Sleep(time.Second * 5)
}
