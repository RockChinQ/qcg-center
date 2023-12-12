package database

// IDatabaseManager 数据库适配器接口
type IDatabaseManager interface {
	// 连接数据库
	Connect() error

	// 储存安装器数据
	StoreInstallerReport(report *InstallerReport) error

	// 储存QChatGPT主程序使用量记录
	StoreQChatGPTUsage(usage *QChatGPTUsage) error

	GetTodayUsageStatic() (TodayUsageStatic, error)
}
