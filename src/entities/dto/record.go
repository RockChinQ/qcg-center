// DTO
package dto

// v2 common
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
	Remote  string `form:"remote" json:"remote" bson:"remote"`
	Author  string `form:"author" json:"author" bson:"author" binding:"required"`
	Version string `form:"version" json:"version" bson:"version" binding:"required"`
}

type RecordDTO struct {
}

// v2 dto
type MainUpdateDTO struct {
	Basic      BasicInfo `form:"basic" json:"basic" bson:"basic" binding:"required"`
	UpdateInfo struct {
		SpentSeconds int    `form:"spent_seconds" json:"spent_seconds" bson:"spent_seconds" binding:"required"`
		InferReason  string `form:"infer_reason" json:"infer_reason" bson:"infer_reason" binding:"required"`
		OldVersion   string `form:"old_version" json:"old_version" bson:"old_version" binding:"required"`
		NewVersion   string `form:"new_version" json:"new_version" bson:"new_version" binding:"required"`
	} `form:"update_info" json:"update_info" bson:"update_info" binding:"required"`
}

type MainAnnouncementDTO struct {
	Basic            BasicInfo `form:"basic" json:"basic" bson:"basic" binding:"required"`
	AnnouncementInfo struct {
		IDs []int `form:"ids" json:"ids" bson:"ids" binding:"required"`
	} `form:"announcement_info" json:"announcement_info" bson:"announcement_info" binding:"required"`
}

type UsageQueryDTO struct {
	Basic       BasicInfo   `form:"basic" json:"basic" bson:"basic" binding:"required"`
	Runtime     RuntimeInfo `form:"runtime" json:"runtime" bson:"runtime" binding:"required"`
	SessionInfo struct {
		Type string `form:"type" json:"type" bson:"type" binding:"required"`
		ID   string `form:"id" json:"id" bson:"id" binding:"required"`
	} `form:"session_info" json:"session_info" bson:"session_info" binding:"required"`
	QueryInfo struct {
		AbilityProvider string `form:"ability_provider" json:"ability_provider" bson:"ability_provider" binding:"required"`
		Usage           int    `form:"usage" json:"usage" bson:"usage" binding:"required"`
		ModelName       string `form:"model_name" json:"model_name" bson:"model_name" binding:"required"`
		ResponseSeconds int    `form:"response_seconds" json:"response_seconds" bson:"response_seconds" binding:"required"`
		RetryTimes      int    `form:"retry_times" json:"retry_times" bson:"retry_times" binding:"required"`
	} `form:"query_info" json:"query_info" bson:"query_info" binding:"required"`
}

type UsageEventDTO struct {
	Basic     BasicInfo    `form:"basic" json:"basic" bson:"basic" binding:"required"`
	Plugins   []PluginInfo `form:"plugins" json:"plugins" bson:"plugins" binding:"required"`
	EventInfo struct {
		Name string `form:"name" json:"name" bson:"name" binding:"required"`
	} `form:"event_info" json:"event_info" bson:"event_info" binding:"required"`
}

type UsageFunctionDTO struct {
	Basic        BasicInfo  `form:"basic" json:"basic" bson:"basic" binding:"required"`
	Plugin       PluginInfo `form:"plugin" json:"plugin" bson:"plugin" binding:"required"`
	FunctionInfo struct {
		Name       string `form:"name" json:"name" bson:"name" binding:"required"`
		Descrition string `form:"description" json:"description" bson:"description" binding:"required"`
	} `form:"function_info" json:"function_info" bson:"function_info" binding:"required"`
}

type PluginInstallDTO struct {
	Basic  BasicInfo  `form:"basic" json:"basic" bson:"basic" binding:"required"`
	Plugin PluginInfo `form:"plugin" json:"plugin" bson:"plugin" binding:"required"`
}

type PluginRemoveDTO struct {
	Basic  BasicInfo  `form:"basic" json:"basic" bson:"basic" binding:"required"`
	Plugin PluginInfo `form:"plugin" json:"plugin" bson:"plugin" binding:"required"`
}

type PluginUpdateDTO struct {
	Basic      BasicInfo  `form:"basic" json:"basic" bson:"basic" binding:"required"`
	Plugin     PluginInfo `form:"plugin" json:"plugin" bson:"plugin" binding:"required"`
	UpdateInfo struct {
		OldVersion string `form:"old_version" json:"old_version" bson:"old_version" binding:"required"`
		NewVersion string `form:"new_version" json:"new_version" bson:"new_version" binding:"required"`
	} `form:"update_info" json:"update_info" bson:"update_info" binding:"required"`
}
