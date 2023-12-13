package v2

import (
	"qcg-center/src/database"
	"strings"

	"github.com/gin-gonic/gin"
)

func Result(c *gin.Context, data interface{}, err error) {
	if err != nil {
		c.JSON(500, gin.H{
			"code": -1,
			"msg":  err.Error(),
		})
	} else {
		if data == nil {
			data = gin.H{}
		}
		c.JSON(200, gin.H{
			"code": 0,
			"msg":  "ok",
			"data": data,
		})
	}
}

func handleRequest[T any](c *gin.Context, db *database.IDatabaseManager, createRecordFunc func(remote_addr string, data *T) error) {
	var report T

	if err := c.ShouldBind(&report); err != nil {
		Result(c, nil, err)
		return
	}
	remoteAddr := c.Request.Header.Get("x-forwarded-for")

	if remoteAddr == "" {
		// 从remoteAddr取
		remoteAddr := c.Request.RemoteAddr
		// 分割IP和端口
		remoteAddrSlice := strings.Split(remoteAddr, ":")
		// 只取IP
		remoteAddr = remoteAddrSlice[0]
	}

	err := createRecordFunc(remoteAddr, &report)
	if err != nil {
		Result(c, nil, err)
		return
	}

	Result(c, nil, nil)
}
