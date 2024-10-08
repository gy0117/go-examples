package main

import (
	"fmt"
	"math/rand"
)

func main() {
	l := NewLru(5)

	// 1. 随机Put三组k-v
	l.randPut(3)
	l.Print()

	// 2. 输出Get(0)
	val := l.Get(0)
	fmt.Printf("Get(0)的结果： %v\n", val)

	// 3. 随机Put两组k-v
	l.randPut(2)
	l.Print()

	// 4. 再Put新的k-v
	l.randPut(1)
	l.Print()

	// 5. Get已经存在的k
	l.Get(0)
	l.Print()

}

func (l *Lru) randPut(num int) {
	for i := 0; i < num; i++ {
		k, v := rand.Intn(20), rand.Intn(50)
		fmt.Printf("randPut, %d - %d\n", k, v)
		l.Put(k, v)
	}
}

type Node struct {
	Key        int
	Val        any
	Prev, Next *Node
}

type Lru struct {
	cache      map[int]*Node
	head, tail *Node
	cap        int
	size       int
}

func NewLru(cap int) *Lru {
	l := &Lru{
		cache: make(map[int]*Node),
		head:  new(Node),
		tail:  new(Node),
		cap:   cap,
	}
	l.head.Next = l.tail
	l.tail.Prev = l.head
	return l
}

func (l *Lru) Get(key int) any {
	node, ok := l.cache[key]
	if !ok {
		return nil
	}
	remove(node)
	l.addToHead(node)
	return node.Val
}

func (l *Lru) Put(key int, val any) {
	node, ok := l.cache[key]
	if ok {
		node.Val = val
		l.cache[key] = node
		remove(node)
		l.addToHead(node)
		return
	}
	newNode := &Node{
		Key: key,
		Val: val,
	}
	l.cache[key] = newNode
	l.size++
	l.addToHead(newNode)
	if l.size > l.cap {
		lastNode := l.tail.Prev
		remove(lastNode)
		delete(l.cache, lastNode.Key)
		l.size--
	}
}

func (l *Lru) Print() {
	cur := l.head.Next
	for cur != l.tail {
		fmt.Printf("%d - %d\n", cur.Key, cur.Val)
		cur = cur.Next
	}
}

func remove(node *Node) {
	node.Prev.Next = node.Next
	node.Next.Prev = node.Prev
}

func (l *Lru) addToHead(node *Node) {
	node.Next = l.head.Next
	l.head.Next.Prev = node

	l.head.Next = node
	node.Prev = l.head
}
