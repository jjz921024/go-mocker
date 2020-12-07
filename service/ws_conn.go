package service

import (
	"errors"
	"github.com/gorilla/websocket"
	"log"
	"sync"
)

type WsConnection struct {
	conn      *websocket.Conn
	inChan    chan []byte // 从通道里读
	outChan   chan []byte // 往通道里写
	closeChan chan byte
	mutex     sync.Mutex
	isClose   bool
}

var ConnSet = make(map[*WsConnection]bool)

func CreateWsConnection(conn *websocket.Conn) (err error) {
	wsConn := &WsConnection{
		conn:      conn,
		inChan:    make(chan []byte, 1000),
		outChan:   make(chan []byte, 1000),
		closeChan: make(chan byte, 1), // 用于感知连接被关闭
	}

	// 启动协程不断从channel里读数据写入websocket
	go func() {
		for true {
			select {
			case data := <-wsConn.outChan:
				if err := wsConn.conn.WriteMessage(websocket.TextMessage, data); err != nil {
					wsConn.Close()
					log.Printf("websocket: %v 写数据异常\n", wsConn.conn)
				}
			case <-wsConn.closeChan:
				wsConn.Close()
			}
		}
	}()

	// 读协程
	go func() {
		for true {
			msgType, data, err := wsConn.conn.ReadMessage()
			if err != nil || msgType == websocket.CloseMessage {
				wsConn.Close()
				log.Printf("websocket连接关闭\n")
				return
			}

			// 读到数据，放入channel
			select {
			case wsConn.inChan <- data:
				log.Printf("ws %v: %s", wsConn, string(data))
			case <-wsConn.closeChan:
				// closeChan被关闭时进入该分支
				wsConn.Close()
			}

		}
	}()

	// 将创建的连接保存到集合中
	ConnSet[wsConn] = true
	return nil
}

func (wsConn *WsConnection) ReadMessage() (data []byte, err error) {
	select {
	case data = <-wsConn.inChan:
	case <-wsConn.closeChan:
		// closeChan可读时，表示连接被关闭
		err = errors.New("connection is closed")
	}
	return data, err
}

func (wsConn *WsConnection) WriteMessage(data []byte) (err error) {
	select {
	case wsConn.outChan <- data:
	case <-wsConn.closeChan:
		err = errors.New("connection is closed")
	}
	return err
}

func (wsConn *WsConnection) Close() {
	// ws的close方法是线程安全的
	_ = wsConn.conn.Close()
	delete(ConnSet, wsConn)

	// 关闭channel，让inChan，outChan在连接关闭时得到信号
	wsConn.mutex.Lock()
	if !wsConn.isClose {
		// 保证通道只被关闭一次
		close(wsConn.closeChan)
		wsConn.isClose = true
	}
	wsConn.mutex.Unlock()
}
