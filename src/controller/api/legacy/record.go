package legacy

import (
	"qcg-center/src/entities/dto"
	"qcg-center/src/service/record"

	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
)

func BindPath(r *gin.Engine, sv *record.RecordService) {

	r.GET("/legacy/usage", func(ctx *gin.Context) {
		var legacyUsageDTO dto.LegacyUsageDTO

		if err := ctx.ShouldBind(&legacyUsageDTO); err != nil {
			ctx.JSON(400, gin.H{"error": err.Error()})
			return
		}

		rid := uuid.New().String()

		usageQueryDTO := dto.UsageQueryDTO{
			Basic: dto.BasicInfo{
				RID:             rid,
				InstanceID:      "from-legacy",
				HostID:          "from-legacy",
				SemanticVersion: legacyUsageDTO.Version,
				Platform:        "from-legacy",
			},
			Runtime: dto.RuntimeInfo{
				AccountID: "from-legacy",
				AdminID:   "from-legacy",
				MsgSource: legacyUsageDTO.MsgSource,
			},
			SessionInfo: struct {
				Type string `form:"type" json:"type" bson:"type" binding:"required"`
				ID   string `form:"id" json:"id" bson:"id" binding:"required"`
			}{
				Type: "from-legacy",
				ID:   "from-legacy",
			},
			QueryInfo: struct {
				AbilityProvider string `form:"ability_provider" json:"ability_provider" bson:"ability_provider" binding:"required"`
				Usage           int    `form:"usage" json:"usage" bson:"usage" binding:"required"`
				ModelName       string `form:"model_name" json:"model_name" bson:"model_name" binding:"required"`
				ResponseSeconds int    `form:"response_seconds" json:"response_seconds" bson:"response_seconds" binding:"required"`
				RetryTimes      int    `form:"retry_times" json:"retry_times" bson:"retry_times" binding:"required"`
			}{
				AbilityProvider: legacyUsageDTO.ServiceName,
				Usage:           legacyUsageDTO.Count,
				ModelName:       "from-legacy",
				ResponseSeconds: -1,
				RetryTimes:      -1,
			},
		}

		if err := sv.InsertRecord(ctx, usageQueryDTO); err != nil {
			ctx.JSON(500, gin.H{"error": err.Error()})
			return
		}

		// 返回纯文本：ok
		ctx.String(200, "ok")
	})
}
