package main

import (
	"fmt"
	"golang.org/x/sync/singleflight"
	"time"
)

func getData(num int) (string, error) {
	fmt.Println("getData")
	time.Sleep(time.Second * 3)
	return fmt.Sprintf("data_%d", num), nil
}

// getData只执行了一次
func main() {
	sg := new(singleflight.Group)

	go func() {
		// shared：为true 表示多个调用共享了返回结果
		v, _, shared := sg.Do("getData", func() (interface{}, error) {
			return getData(1)
		})
		fmt.Printf("one call v: %s, shared: %v\n", v.(string), shared)
	}()

	time.Sleep(time.Millisecond * 500)

	go func() {
		v, _, shared := sg.Do("getData", func() (interface{}, error) {
			return getData(2)
		})
		fmt.Printf("two call v: %s, shared: %v\n", v.(string), shared)
	}()

	time.Sleep(time.Minute)
}
