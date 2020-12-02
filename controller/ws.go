package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go_mocker/service"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	// 允许跨域
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// 处理websocket升级
func WsUpgradeHandler(c *gin.Context) {
	var conn, err = upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("websocket握手失败: %s\n", err.Error())
		return
	}

	// todo 断开连接触发哪个函数
	/*conn.SetCloseHandler(func(code int, text string) error {
		println("close....")
		_ = conn.Close()
		return nil
	  })*/

	err = service.CreateWsConnection(conn)
	if err != nil {
		log.Printf("websocket握手失败: %s\n", err.Error())
	}
}
