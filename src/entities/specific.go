package entities

type MainUpdate struct {
	Basic      BasicInfo `form:"basic" json:"basic" bson:"basic" binding:"required"`
	UpdateInfo struct {
		SpentSeconds int `form:"spent_seconds" json:"spent_seconds" bson:"spent_seconds" binding:"required"`
		InferReason  int `form:"infer_reason" json:"infer_reason" bson:"infer_reason" binding:"required"`
		OldVersion   int `form:"old_version" json:"old_version" bson:"old_version" binding:"required"`
		NewVersion   int `form:"new_version" json:"new_version" bson:"new_version" binding:"required"`
	} `form:"update_info" json:"update_info" bson:"update_info" binding:"required"`
}

type MainAnnouncement struct {
	Basic            BasicInfo `form:"basic" json:"basic" bson:"basic" binding:"required"`
	AnnouncementInfo struct {
		IDs []int `form:"ids" json:"ids" bson:"ids" binding:"required"`
	} `form:"announcement_info" json:"announcement_info" bson:"announcement_info" binding:"required"`
}

type UsageQuery struct {
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

type UsageEvent struct {
	Basic     BasicInfo    `form:"basic" json:"basic" bson:"basic" binding:"required"`
	Plugins   []PluginInfo `form:"plugins" json:"plugins" bson:"plugins" binding:"required"`
	EventInfo struct {
		Name string `form:"name" json:"name" bson:"name" binding:"required"`
	} `form:"event_info" json:"event_info" bson:"event_info" binding:"required"`
}

type UsageFunction struct {
	Basic        BasicInfo  `form:"basic" json:"basic" bson:"basic" binding:"required"`
	Plugin       PluginInfo `form:"plugin" json:"plugin" bson:"plugin" binding:"required"`
	FunctionInfo struct {
		Name       string `form:"name" json:"name" bson:"name" binding:"required"`
		Descrition string `form:"description" json:"description" bson:"description" binding:"required"`
	} `form:"function_info" json:"function_info" bson:"function_info" binding:"required"`
}

type PluginInstall struct {
	Basic  BasicInfo  `form:"basic" json:"basic" bson:"basic" binding:"required"`
	Plugin PluginInfo `form:"plugin" json:"plugin" bson:"plugin" binding:"required"`
}

type PluginRemove struct {
	Basic  BasicInfo  `form:"basic" json:"basic" bson:"basic" binding:"required"`
	Plugin PluginInfo `form:"plugin" json:"plugin" bson:"plugin" binding:"required"`
}

type PluginUpdate struct {
	Basic      BasicInfo  `form:"basic" json:"basic" bson:"basic" binding:"required"`
	Plugin     PluginInfo `form:"plugin" json:"plugin" bson:"plugin" binding:"required"`
	UpdateInfo struct {
		OldVersion int `form:"old_version" json:"old_version" bson:"old_version" binding:"required"`
		NewVersion int `form:"new_version" json:"new_version" bson:"new_version" binding:"required"`
	} `form:"update_info" json:"update_info" bson:"update_info" binding:"required"`
}
