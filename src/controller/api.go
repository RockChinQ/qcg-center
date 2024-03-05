package controller

import (
	"log"
	database "qcg-center/src/database"
	util "qcg-center/src/util"
	"strconv"
	"strings"

	legacy "qcg-center/src/controller/legacy"
	v2 "qcg-center/src/controller/v2"

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
		var report legacy.LegacyReport

		if err := c.ShouldBind(&report); err != nil {
			log.Println(err)
			c.JSON(400, gin.H{"error": "missing params."})
			return
		}

		// 保存到数据库
		var installerReport database.InstallerReport

		installerReport.OSName = report.OSName
		installerReport.Arch = report.Arch

		installerReport.Timestamp = util.GetCSTTime().Unix()
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
		var usage legacy.LegacyUsage

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

		qchatgptUsage.Timestamp = util.GetCSTTime().Unix()

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

	mainUpdate := v2.MainUpdate(&m.dbmgr)
	mainAnnouncement := v2.MainAnnouncement(&m.dbmgr)

	usageQuery := v2.UsageQuery(&m.dbmgr)
	usageEvent := v2.UsageEvent(&m.dbmgr)
	usageFunction := v2.UsageFunction(&m.dbmgr)

	pluginInstall := v2.PluginInstall(&m.dbmgr)
	pluginRemove := v2.PluginRemove(&m.dbmgr)
	pluginUpdate := v2.PluginUpdate(&m.dbmgr)

	r.POST("/api/v2/main/update", mainUpdate)
	r.POST("/api/v2/main/announcement", mainAnnouncement)

	r.POST("/api/v2/usage/query", usageQuery)
	r.POST("/api/v2/usage/event", usageEvent)
	r.POST("/api/v2/usage/function", usageFunction)

	r.POST("/api/v2/plugin/install", pluginInstall)
	r.POST("/api/v2/plugin/remove", pluginRemove)
	r.POST("/api/v2/plugin/update", pluginUpdate)

	grafanaRoot := legacy.GrafanaRoot(&m.dbmgr)
	grafanaTodayUsageStatic := legacy.GrafanaTodayUsageStatic(&m.dbmgr)
	grafanaRecentDaysUsageTrend := legacy.GrafanaRecentDaysUsageTrend(&m.dbmgr)

	r.GET("/grafana", grafanaRoot)
	r.POST("/grafana", grafanaRoot)

	r.GET("/grafana/today_usage_static", grafanaTodayUsageStatic)
	r.POST("/grafana/today_usage_static", grafanaTodayUsageStatic)

	r.GET("/grafana/recent_days_usage_trend", grafanaRecentDaysUsageTrend)
	r.POST("/grafana/recent_days_usage_trend", grafanaRecentDaysUsageTrend)

	m.r = r

	return nil
}

// 启动WebAPI服务
func (m *WebAPI) Serve() error {
	log.Println("WebAPI listening on " + m.addr + ":" + strconv.Itoa(m.port))
	return m.r.Run(m.addr + ":" + strconv.Itoa(m.port))
}
