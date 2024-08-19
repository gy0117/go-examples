package main

import (
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"sync"
)

// websocket的read, write message都不是线程安全的
// websocket的close是线程安全的

// 1. 封装wsConn
// 2. 从inChan读消息，发送消息到outChan（对于调用者来说的）

type wsMsg struct {
	msgType int
	data    []byte
}

type Connection struct {
	wsConn  *websocket.Conn
	inChan  chan *wsMsg
	outChan chan *wsMsg

	mux       sync.Mutex
	isClosed  bool
	closeChan chan struct{}
}

func NewConnection(wsConn *websocket.Conn) *Connection {
	conn := &Connection{
		wsConn:    wsConn,
		inChan:    make(chan *wsMsg, 1000),
		outChan:   make(chan *wsMsg, 1000),
		closeChan: make(chan struct{}, 1),
	}

	go conn.readLoop()
	go conn.writeLoop()

	return conn
}

func (conn *Connection) Close() {
	// websocket的close是线程安全的，可冲入
	conn.wsConn.Close()

	// 重复执行会报错，只能执行一次
	conn.mux.Lock()
	defer conn.mux.Unlock()
	if !conn.isClosed {
		close(conn.closeChan)
		conn.isClosed = true
	}
}

// API方法
func (conn *Connection) ReadMessage() (data []byte, err error) {
	select {
	case msg := <-conn.inChan:
		data = msg.data
	case <-conn.closeChan:
		err = errors.New("connection is closed")
	}
	return
}

// API方法，需要做到线程安全，因此使用chan，读消息同理
func (conn *Connection) WriteMessage(data []byte) (err error) {
	select {
	case conn.outChan <- &wsMsg{
		msgType: websocket.TextMessage,
		data:    data,
	}:
	case <-conn.closeChan:
		err = errors.New("connection is closed")
	}
	return
}

func (conn *Connection) readLoop() {
	for {
		// 从websocket里读消息
		msgType, data, err := conn.wsConn.ReadMessage()
		if err != nil {
			conn.Close()
			return
		}
		msg := &wsMsg{
			msgType: msgType,
			data:    data,
		}
		// 假设这个写满了，就会阻塞在这里，然后WriteMessage时报错，整个conn被close了
		// 但是此时这个不知道情况，仍然阻塞着
		//conn.inChan <- msg

		select {
		case conn.inChan <- msg:
		case <-conn.closeChan:
			// 当closeChan被关闭时，进入这个分支
			conn.Close()
			return
		}
	}
}

func (conn *Connection) writeLoop() {
	for {
		var msg *wsMsg
		select {
		case msg = <-conn.outChan:
		case <-conn.closeChan:
			conn.Close()
			return
		}
		err := conn.wsConn.WriteMessage(msg.msgType, msg.data)
		if err != nil {
			conn.Close()
			return
		}
	}
}
