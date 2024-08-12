package main

import (
	"context"
	"crypto/sha256"
	"log"
	"math/big"
	"os"
	"strconv"
	"time"
)

// Context主动撤销
// pow执行时间过长，主goroutine撤销子goroutine的执行
func main() {
	targetBits, _ := strconv.Atoi(os.Args[1])
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)

	pow := func(ctx context.Context, targetBits int, ch chan string) {
		target := big.NewInt(1)
		// 除了前targetBits位，其余位都是1
		target.Lsh(target, uint(256-targetBits))

		var hashInt big.Int
		var hash [32]byte
		nonce := 0

		// 寻找一个满足当前难度的数
		for {
			select {
			case <-ctx.Done():
				log.Println("context is Canceled")
				ch <- ""
				return
			default:
				data := "hello world" + strconv.Itoa(nonce)
				hash = sha256.Sum256(hash[:]) // 计算hash值
				hashInt.SetBytes(hash[:])     // 将hash值转换位bit.Int

				// hashInt <= target，找到一个不大于目标值的数，前targetBits位都为0
				if hashInt.Cmp(target) < 1 {
					ch <- data
					return
				} else {
					nonce++
				}
			}
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch := make(chan string, 1)

	log.Println("开始寻找一个数，使得hash值小于目标值")
	go pow(ctx, targetBits, ch)

	time.Sleep(time.Second)

	// select会等待任意一个非阻塞的通道操作。
	// default子句的存在使得select语句变得非阻塞，也就是说如果所有的通道操作都阻塞，select不会等待，而是直接执行default中的代码。
	// 如果没有default子句且所有通道操作都阻塞，select语句将会阻塞，直到至少有一个通道操作可以执行。
	select {
	case result := <-ch:
		log.Println("找到一个比目标值小的数： ", result)
		return
	default:
		log.Println("没有找到比目标值小的数: ", ctx.Err())
	}
}
