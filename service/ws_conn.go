package service

import (
	"github.com/gorilla/websocket"
	"log"
)

type WsConnection struct {
	conn    *websocket.Conn
	inChan  chan []byte // 从通道里读
	outChan chan []byte // 往通道里写
}

var ConnSet = make(map[*WsConnection]bool)

func CreateWsConnection(conn *websocket.Conn) (err error) {
	wsConn := &WsConnection{
		conn:    conn,
		inChan:  make(chan []byte, 1000),
		outChan: make(chan []byte, 1000),
	}

	// 启动协程不断从channel里读数据写入websocket
	go func() {
		for true {
			var data = <-wsConn.outChan
			if err := wsConn.conn.WriteMessage(websocket.TextMessage, data); err != nil {
				wsConn.Close()
				log.Printf("websocket: %v 写数据异常\n", wsConn.conn)
			}
		}
	}()

	go func() {
		for true {
			msgType, _, err := wsConn.conn.ReadMessage()
			if err != nil || msgType == websocket.CloseMessage {
				wsConn.Close()
				log.Printf("websocket连接关闭\n")
				return
			}

			//wsConn.inChan <- data
			//log.Printf("ws %v: %s", wsConn, string(data))
		}
	}()

	// 将创建的连接保存到集合中
	ConnSet[wsConn] = true
	return nil
}

func (wsConn *WsConnection) ReadMessage() (data []byte, err error) {
	data = <-wsConn.inChan
	return data, nil
}

func (wsConn *WsConnection) WriteMessage(data []byte) (err error) {
	wsConn.outChan <- data
	return nil
}

func (wsConn *WsConnection) Close() {
	// ws的close方法是线程安全的
	_ = wsConn.conn.Close()
	delete(ConnSet, wsConn)
}
