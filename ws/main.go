package main

import (
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

var upgrader = websocket.Upgrader{
	// 允许跨域，即跨域名
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// ws demo
func main() {

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		//var msgType  int
		//var data []byte
		var err error
		// 回调函数
		//w.Write([]byte("Hello 0xAAC"))
		// 升级，就是往header里面加 Upgrade: websocket等数据
		// 返回的conn，就是websocket长链接
		wsConn, _ := upgrader.Upgrade(w, r, nil)
		conn := NewConnection(wsConn)

		//for {
		//	_, data, err = conn.ReadMessage()
		//	if err != nil {
		//		conn.Close()
		//		return
		//	}
		//	if err = conn.WriteMessage(websocket.TextMessage, data); err != nil {
		//		conn.Close()
		//		return
		//	}
		//}

		go func() {
			// 测试一下线程安全
			for {
				err = conn.WriteMessage([]byte("heartbeat"))
				if err != nil {
					return
				}
				time.Sleep(time.Second * 2)
			}

		}()

		for {
			msg, err := conn.ReadMessage()
			if err != nil {
				conn.Close()
				return
			}
			// 将数据返回回去
			err = conn.WriteMessage(msg)
			if err != nil {
				conn.Close()
				return
			}
		}

	})

	_ = http.ListenAndServe("127.0.0.1:7777", nil)

}
