package record

import (
	"qcg-center/src/service/record"

	"github.com/gin-gonic/gin"

	"qcg-center/src/entities/dto"
)

func CommonRecordHandlerGeneric[T any](sv *record.RecordService) func(c *gin.Context) {
	return func(c *gin.Context) {

		var reportDTO T

		if err := c.ShouldBindJSON(&reportDTO); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		if err := sv.InsertRecord(c, reportDTO); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{
			"code": 0,
			"msg":  "ok",
			"data": gin.H{},
		})
	}
}

func BindPath(r *gin.Engine, sv *record.RecordService) {
	// 方法、路径、DTO映射
	r.POST("/api/v2/main/update", CommonRecordHandlerGeneric[dto.MainUpdateDTO](sv))
	r.POST("/api/v2/main/announcement", CommonRecordHandlerGeneric[dto.MainAnnouncementDTO](sv))

	r.POST("/api/v2/usage/query", CommonRecordHandlerGeneric[dto.UsageQueryDTO](sv))
	r.POST("/api/v2/usage/event", CommonRecordHandlerGeneric[dto.UsageEventDTO](sv))
	r.POST("/api/v2/usage/function", CommonRecordHandlerGeneric[dto.UsageFunctionDTO](sv))

	r.POST("/api/v2/plugin/install", CommonRecordHandlerGeneric[dto.PluginInstallDTO](sv))
	r.POST("/api/v2/plugin/remove", CommonRecordHandlerGeneric[dto.PluginRemoveDTO](sv))
	r.POST("/api/v2/plugin/update", CommonRecordHandlerGeneric[dto.PluginUpdateDTO](sv))
}
