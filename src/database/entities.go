// 数据库层面数据结构
package database

import "time"

// 安装器报告数据
type InstallerReport struct {
	OSName     string `bson:"osname"`
	Arch       string `bson:"arch"`
	Timestamp  int64  `bson:"timestamp"`
	Mac        string `bson:"mac"`
	Version    string `bson:"version"`
	Message    string `bson:"message"`
	RemoteAddr string `bson:"remote_addr"`
}

// QChatGPT主程序使用量记录
type QChatGPTUsage struct {
	ServiceName string `bson:"service_name"`
	Version     string `bson:"version"`
	Count       int    `bson:"count"`
	MsgSource   string `bson:"msg_source"`
	Timestamp   int64  `bson:"timestamp"`
	RemoteAddr  string `bson:"remote_addr"`
}

type DailyUsageStatic struct {
	Number          int       `bson:"number" json:"number"`
	Begin           time.Time `bson:"begin" json:"begin"`
	Duration        int64     `bson:"duration" json:"duration"`
	UsageCount      int       `bson:"usage_count" json:"usage_count"`
	ActiveHostCount int       `bson:"active_host_count" json:"active_host_count"`
	NewHostCount    int       `bson:"new_host_count" json:"new_host_count"`
	ModifiedAt      time.Time `bson:"modified_at" json:"modified_at"`
}

type CommonDocumnet struct {
	RemoteAddr string      `bson:"remote_addr" json:"remote_addr"`
	Time       time.Time   `bson:"time" json:"time"`
	Data       interface{} `bson:"data" json:"data"`
}

type HostIDRecord struct {
	HostID string    `bson:"host_id" json:"host_id"`
	Time   time.Time `bson:"time" json:"time"`
}

type InstanceIDRecord struct {
	InstanceID string    `bson:"instance_id" json:"instance_id"`
	Time       time.Time `bson:"time" json:"time"`
}

type IPRecord struct {
	IP   string    `bson:"ip" json:"ip"`
	Time time.Time `bson:"time" json:"time"`
}

type HostInstanceIPTuple struct {
	HostID     string    `bson:"host_id" json:"host_id"`
	InstanceID string    `bson:"instance_id" json:"instance_id"`
	IP         string    `bson:"ip" json:"ip"`
	Time       time.Time `bson:"time" json:"time"`
}
