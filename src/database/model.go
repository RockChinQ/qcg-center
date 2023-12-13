package database

import "qcg-center/src/entities"

// IDatabaseManager 数据库适配器接口
type IDatabaseManager interface {
	// 连接数据库
	Connect() error

	// 储存安装器数据
	StoreInstallerReport(report *InstallerReport) error

	// 储存QChatGPT主程序使用量记录
	StoreQChatGPTUsage(usage *QChatGPTUsage) error

	GetTodayUsageStatic() (DailyUsageStatic, error)

	GetRecentDaysUsageTrend(day int) ([]DailyUsageStatic, error)

	// v2
	InsertMainUpdateRecord(remote_addr string, record *entities.MainUpdate) error
	InsertMainAnnouncementRecord(remote_addr string, record *entities.MainAnnouncement) error

	InsertUsageQueryRecord(remote_addr string, record *entities.UsageQuery) error
	InsertUsageEventRecord(remote_addr string, record *entities.UsageEvent) error
	InsertUsageFunctionRecord(remote_addr string, record *entities.UsageFunction) error

	InsertPluginInstallRecord(remote_addr string, record *entities.PluginInstall) error
	InsertPluginRemoveRecord(remote_addr string, record *entities.PluginRemove) error
	InsertPluginUpdateRecord(remote_addr string, record *entities.PluginUpdate) error
}
