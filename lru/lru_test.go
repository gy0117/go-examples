package main

import (
	"fmt"
	"testing"
)

var l = NewLru(5)

// 1. 随机Put三组k-v
func TestRandPut3(t *testing.T) {
	l.randPut(3)
	l.Print()
}

func TestGet(t *testing.T) {
	//l.randPut(3)
	//l.Print()
	val := l.Get(0)
	fmt.Printf("Get(0)的结果： %v\n", val)
}

func TestRandPut2(t *testing.T) {
	l.randPut(3)
	l.randPut(2)
	l.Print()
}

func TestRandPut(t *testing.T) {
	l.randPut(3)
	l.randPut(2)
	l.randPut(1)
	l.Print()
}

func TestGetHasExist(t *testing.T) {
	l.Put(1, 21)
	val := l.Get(1)
	fmt.Printf("Get的结果： %v\n", val)
}
