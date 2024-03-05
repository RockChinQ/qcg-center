package legacy

type LegacyReport struct {
	OSName    string `form:"osname" json:"osname" binding:"required"`
	Arch      string `form:"arch" json:"arch" binding:"required"`
	Timestamp int64  `form:"timestamp" json:"timestamp" binding:"required"`
	Version   string `form:"version" json:"version" binding:"required"`
	Message   string `form:"message" json:"message"`
}

type LegacyUsage struct {
	ServiceName string `form:"service_name" json:"service_name" binding:"required"`
	Version     string `form:"version" json:"version" binding:"required"`
	Count       int    `form:"count" json:"count" binding:"required"`
	MsgSource   string `form:"msg_source" json:"msg_source"`
}
