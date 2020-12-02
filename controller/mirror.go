package controller

import (
	"github.com/gin-gonic/gin"
	"go_mocker/service"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func Mirror(c *gin.Context) {
	sb := new(strings.Builder)

	sb.WriteString("---------------------[" + time.Now().Format("2006-01-02 15:04:05") + "]---------------------\n")
	sb.WriteString("[流量拷贝] -> " + c.Request.RemoteAddr + "\n")
	sb.WriteString(c.Request.Method + " " + c.Request.RequestURI + " " + c.Request.Proto + "\n")

	sb.WriteString("\n")
	for header := range c.Request.Header {
		sb.WriteString(header + ":" + c.Request.Header.Get(header) + "\n")
	}

	b, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	sb.WriteString("\n")
	sb.Write(b)
	sb.WriteString("\n")

	c.Status(200)

	go func() {
		for conn := range service.ConnSet {
			_ = conn.WriteMessage([]byte(sb.String()))
		}
	}()
}
