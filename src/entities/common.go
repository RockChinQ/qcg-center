package entities

// v2
type BasicInfo struct {
	RID             string `form:"rid" json:"rid" bson:"rid" binding:"required"`
	InstanceID      string `form:"instance_id" json:"instance_id" bson:"instance_id" binding:"required"`
	HostID          string `form:"host_id" json:"host_id" bson:"host_id" binding:"required"`
	SemanticVersion string `form:"semantic_version" json:"semantic_version" bson:"semantic_version" binding:"required"`
	Platform        string `form:"platform" json:"platform" bson:"platform" binding:"required"`
}

type RuntimeInfo struct {
	AccountID string `form:"account_id" json:"account_id" bson:"account_id" binding:"required"`
	AdminID   string `form:"admin_id" json:"admin_id" bson:"admin_id" binding:"required"`
	MsgSource string `form:"msg_source" json:"msg_source" bson:"msg_source" binding:"required"`
}

type PluginInfo struct {
	Name    string `form:"name" json:"name" bson:"name" binding:"required"`
	Remote  string `form:"remote" json:"remote" bson:"remote" binding:"required"`
	Author  string `form:"author" json:"author" bson:"author" binding:"required"`
	Version string `form:"version" json:"version" bson:"version" binding:"required"`
}
