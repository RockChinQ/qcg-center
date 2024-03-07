package api

import (
	"log"
	"strconv"

	"qcg-center/src/controller/api/legacy"
	"qcg-center/src/controller/api/v2/record"
	"qcg-center/src/controller/api/v2/view"
	serviceRecord "qcg-center/src/service/record"
	serviceView "qcg-center/src/service/view"

	"github.com/gin-gonic/gin"
)

type WebAPI struct {

	// 服务
	SvRecord *serviceRecord.RecordService

	SvView *serviceView.RealTimeDataService

	// port
	Port int

	// addr
	Addr string

	R *gin.Engine
}

func NewWebAPI(svRecord *serviceRecord.RecordService, svView *serviceView.RealTimeDataService, port int, addr string) *WebAPI {
	r := gin.Default()

	record.BindPath(r, svRecord)

	view.BindPath(r, svView)

	legacy.BindPath(r, svRecord)

	view_root := func(c *gin.Context) {
		c.JSON(200, gin.H{
			"code": 0,
			"msg":  "ok",
			"data": gin.H{},
		})
	}

	r.GET("/api/v2/view/", view_root)
	r.POST("/api/v2/view/", view_root)

	webapi := &WebAPI{
		SvRecord: svRecord,
		SvView:   svView,
		Port:     port,
		Addr:     addr,
		R:        r,
	}

	return webapi
}

// 启动WebAPI服务
func (m *WebAPI) Serve() error {
	log.Println("WebAPI listening on " + m.Addr + ":" + strconv.Itoa(m.Port))
	return m.R.Run(m.Addr + ":" + strconv.Itoa(m.Port))
}
