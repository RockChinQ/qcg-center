package api

import (
	"log"
	database "qcg-center/src/database"
	util "qcg-center/src/util"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type WebAPI struct {
	Cfg *util.Config

	// 数据库管理器
	dbmgr database.IDatabaseManager

	// port
	port int

	// addr
	addr string

	r *gin.Engine
}

// 初始化WebAPI
func (m *WebAPI) Init(dbmgr database.IDatabaseManager) error {

	m.dbmgr = dbmgr

	m.port = m.Cfg.API.Port
	m.addr = m.Cfg.API.Listen

	r := gin.Default()

	r.GET("/legacy/report", func(c *gin.Context) {
		var report LegacyReport

		if err := c.ShouldBind(&report); err != nil {
			log.Println(err)
			c.JSON(400, gin.H{"error": "missing params."})
			return
		}

		// 保存到数据库
		var installerReport database.InstallerReport

		installerReport.OSName = report.OSName
		installerReport.Arch = report.Arch
		installerReport.Timestamp = report.Timestamp
		installerReport.Version = report.Version
		installerReport.Message = report.Message

		// 从header取到x-forwarded-for
		installerReport.RemoteAddr = c.Request.Header.Get("x-forwarded-for")

		if installerReport.RemoteAddr == "" {
			// installerReport.RemoteAddr = c.Request.RemoteAddr
			remoteAddr := c.Request.RemoteAddr
			// 分割IP和端口
			remoteAddrSlice := strings.Split(remoteAddr, ":")
			// 只取IP
			installerReport.RemoteAddr = remoteAddrSlice[0]
		}

		err := m.dbmgr.StoreInstallerReport(&installerReport)

		if err != nil {
			log.Println(err)
			c.JSON(500, gin.H{"error": "internal error."})
			return
		}

		// code 200
		c.Writer.WriteHeader(200)
		// 返回纯文字ok，不是json
		c.Writer.WriteString("ok")
	})

	r.GET("/legacy/usage", func(c *gin.Context) {
		var usage LegacyUsage

		if err := c.ShouldBind(&usage); err != nil {
			log.Println(err)
			c.JSON(400, gin.H{"error": "missing params."})
			return
		}

		// 保存到数据库
		var qchatgptUsage database.QChatGPTUsage

		qchatgptUsage.ServiceName = usage.ServiceName
		qchatgptUsage.Version = usage.Version
		qchatgptUsage.Count = usage.Count
		qchatgptUsage.MsgSource = usage.MsgSource
		qchatgptUsage.Timestamp = time.Now().Unix()

		// 从header取到x-forwarded-for
		qchatgptUsage.RemoteAddr = c.Request.Header.Get("x-forwarded-for")

		if qchatgptUsage.RemoteAddr == "" {
			// 从remoteAddr取
			remoteAddr := c.Request.RemoteAddr
			// 分割IP和端口
			remoteAddrSlice := strings.Split(remoteAddr, ":")
			// 只取IP
			qchatgptUsage.RemoteAddr = remoteAddrSlice[0]
		}

		err := m.dbmgr.StoreQChatGPTUsage(&qchatgptUsage)

		if err != nil {
			log.Println(err)
			c.JSON(500, gin.H{"error": "internal error."})
			return
		}

		// code 200
		c.Writer.WriteHeader(200)
		// 返回纯文字ok，不是json
		c.Writer.WriteString("ok")
	})

	m.r = r

	return nil
}

// 启动WebAPI服务
func (m *WebAPI) Serve() error {
	log.Println("WebAPI listening on " + m.addr + ":" + strconv.Itoa(m.port))
	return m.r.Run(m.addr + ":" + strconv.Itoa(m.port))
}
