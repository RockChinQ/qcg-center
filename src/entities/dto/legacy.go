package dto

type LegacyUsageDTO struct {
	ServiceName string `form:"service_name" json:"service_name" binding:"required"`
	Version     string `form:"version" json:"version" binding:"required"`
	Count       int    `form:"count" json:"count" binding:"required"`
	MsgSource   string `form:"msg_source" json:"msg_source"`
}
