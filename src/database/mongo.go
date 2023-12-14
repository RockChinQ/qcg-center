package database

import (
	"context"
	"time"

	"qcg-center/src/entities"
	"qcg-center/src/util"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	DIRECT_MAIN_UPDATE_RECORD_COLLECTION_NAME = "direct_main_update_record"
	DIRECT_MAIN_ANNOUNCEMENT_COLLECTION_NAME  = "direct_main_announcement"
	DIRECT_USAGE_QUERY_COLLECTION_NAME        = "direct_usage_query"
	DIRECT_USAGE_EVENT_COLLECTION_NAME        = "direct_usage_event"
	DIRECT_USAGE_FUNCTION_COLLECTION_NAME     = "direct_usage_function"
	DIRECT_PLUGIN_INSTALL_COLLECTION_NAME     = "direct_plugin_install"
	DIRECT_PLUGIN_REMOVE_COLLECTION_NAME      = "direct_plugin_remove"
	DIRECT_PLUGIN_UPDATE_COLLECTION_NAME      = "direct_plugin_update"

	ANALYSIS_HOST_ID_COLLECTION_NAME          = "analysis_host_id"
	ANALYSIS_INSTANCE_ID_COLLECTION_NAME      = "analysis_instance_id"
	ANALYSIS_IP_COLLECTION_NAME               = "analysis_ip"
	ANALYSIS_HOST_INSTANCE_IP_COLLECTION_NAME = "analysis_host_instance_ip"

	VIEW_DAILY_USAGE_REPORT_COLLECTION_NAME = "view_daily_usage_report"
)

type MongoDBManager struct {
	Cfg *util.Config

	Client   *mongo.Client
	Database string
}

func (m *MongoDBManager) Connect() error {
	uri := m.Cfg.Database.Params["uri"]

	m.Database = m.Cfg.Database.Params["database"]

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))

	if err != nil {
		return err
	}

	m.Client = client

	return nil
}

func (m *MongoDBManager) StoreInstallerReport(report *InstallerReport) error {
	coll := m.Client.Database("qcg-center-records").Collection("installer-reports")

	_, err := coll.InsertOne(context.TODO(), report)

	if err != nil {
		return err
	}

	// 从analysis_reports_remote_addrs中查找是否有相同的remote_addr
	// 无则插入(remote_addr, 当前时间)
	check, err := m.Client.Database("qcg-center-records").Collection("analysis_reports_remote_addrs").Find(context.TODO(), map[string]interface{}{
		"remote_addr": report.RemoteAddr,
	})

	if err != nil {
		return err
	}

	if !check.Next(context.Background()) {
		// 无记录
		_, err = m.Client.Database("qcg-center-records").Collection("analysis_reports_remote_addrs").InsertOne(context.TODO(), map[string]interface{}{
			"remote_addr": report.RemoteAddr,
			// 使用report.Timestamp
			"created_at": util.GetCSTTime(),
		})
	}

	return err
}

func (m *MongoDBManager) StoreQChatGPTUsage(usage *QChatGPTUsage) error {
	coll := m.Client.Database("qcg-center-records").Collection("qchatgpt-usage")

	_, err := coll.InsertOne(context.TODO(), usage)

	if err != nil {
		return err
	}

	// 从analysis_usage_remote_addrs中查找是否有相同的remote_addr
	// 无则插入(remote_addr, 当前时间)
	check, err := m.Client.Database("qcg-center-records").Collection("analysis_usage_remote_addrs").Find(context.TODO(), map[string]interface{}{
		"remote_addr": usage.RemoteAddr,
	})

	if err != nil {
		return err
	}

	if !check.Next(context.Background()) {
		// 无记录
		_, err = m.Client.Database("qcg-center-records").Collection("analysis_usage_remote_addrs").InsertOne(context.TODO(), map[string]interface{}{
			"remote_addr": usage.RemoteAddr,
			"created_at":  util.GetCSTTime(),
		})
	}

	return err
}

func (m *MongoDBManager) GetTodayUsageStatic() (DailyUsageStatic, error) {
	var todayUsageStatic DailyUsageStatic

	coll := m.Client.Database("qcg-center-records").Collection("analysis_daily")

	// 统一 UTC
	today := time.Now().UTC()

	// 今天的0点
	today = time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, util.GetCSTTimeLocation())

	res := coll.FindOne(context.TODO(), map[string]interface{}{
		"begin": today,
	})

	if res.Err() != nil {
		return todayUsageStatic, res.Err()
	}

	err := res.Decode(&todayUsageStatic)

	if err != nil {
		return todayUsageStatic, err
	}

	return todayUsageStatic, nil
}

func (m *MongoDBManager) GetRecentDaysUsageTrend(day int) ([]DailyUsageStatic, error) {
	var trend []DailyUsageStatic

	coll := m.Client.Database("qcg-center-records").Collection("analysis_daily")

	// 统一 UTC
	today := time.Now().UTC()

	// 今天的0点
	today = time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, util.GetCSTTimeLocation())

	seq := 0

	// 从今天开始往前推day天
	for i := 0; i < day; i++ {
		res := coll.FindOne(context.TODO(), map[string]interface{}{
			"begin": today.Add(-24 * time.Hour * time.Duration(i)),
		})

		if res.Err() != nil {
			return trend, res.Err()
		}

		var analysis DailyUsageStatic

		err := res.Decode(&analysis)

		if err != nil {
			return trend, err
		}

		analysis.Number = seq
		seq++

		trend = append(trend, analysis)
	}

	return trend, nil
}

func (m *MongoDBManager) saveHostInstanceIPTuple(host_id string, instance_id string, ip string) error {
	// 检查每个表中是否有相同的记录
	// 无则插入(host_id, 当前时间)
	check, err := m.Client.Database(m.Database).Collection(ANALYSIS_HOST_ID_COLLECTION_NAME).Find(context.TODO(), map[string]interface{}{
		"host_id": host_id,
	})

	if err != nil {
		return err
	}

	if !check.Next(context.Background()) {
		// 无记录
		_, err = m.Client.Database(m.Database).Collection(ANALYSIS_HOST_ID_COLLECTION_NAME).InsertOne(context.TODO(), map[string]interface{}{
			"host_id":    host_id,
			"created_at": util.GetCSTTime(),
		})

		if err != nil {
			return err
		}
	}

	check, err = m.Client.Database(m.Database).Collection(ANALYSIS_INSTANCE_ID_COLLECTION_NAME).Find(context.TODO(), map[string]interface{}{
		"instance_id": instance_id,
	})

	if err != nil {
		return err
	}

	if !check.Next(context.Background()) {
		// 无记录
		_, err = m.Client.Database(m.Database).Collection(ANALYSIS_INSTANCE_ID_COLLECTION_NAME).InsertOne(context.TODO(), map[string]interface{}{
			"instance_id": instance_id,
			"created_at":  util.GetCSTTime(),
		})

		if err != nil {
			return err
		}
	}

	check, err = m.Client.Database(m.Database).Collection(ANALYSIS_IP_COLLECTION_NAME).Find(context.TODO(), map[string]interface{}{
		"ip": ip,
	})

	if err != nil {
		return err
	}

	if !check.Next(context.Background()) {
		// 无记录
		_, err = m.Client.Database(m.Database).Collection(ANALYSIS_IP_COLLECTION_NAME).InsertOne(context.TODO(), map[string]interface{}{
			"ip":         ip,
			"created_at": util.GetCSTTime(),
		})

		if err != nil {
			return err
		}
	}

	check, err = m.Client.Database(m.Database).Collection(ANALYSIS_HOST_INSTANCE_IP_COLLECTION_NAME).Find(context.TODO(), map[string]interface{}{
		"host_id":     host_id,
		"instance_id": instance_id,
		"ip":          ip,
	})

	if err != nil {
		return err
	}

	if !check.Next(context.Background()) {
		// 无记录
		_, err = m.Client.Database(m.Database).Collection(ANALYSIS_HOST_INSTANCE_IP_COLLECTION_NAME).InsertOne(context.TODO(), map[string]interface{}{
			"host_id":     host_id,
			"instance_id": instance_id,
			"ip":          ip,
			"created_at":  util.GetCSTTime(),
		})
	}

	return err
}

func (m *MongoDBManager) handleV2DirectData(remote_addr string, basic entities.BasicInfo, record interface{}, coll string) error {
	m.saveHostInstanceIPTuple(basic.HostID, basic.InstanceID, remote_addr)
	// 保存到对应的collection中
	_, err := m.Client.Database(m.Database).Collection(coll).InsertOne(context.TODO(), map[string]interface{}{
		"remote_addr": remote_addr,
		"time":        util.GetCSTTime(),
		"data":        record,
	})

	return err
}

// v2
func (m *MongoDBManager) InsertMainUpdateRecord(remote_addr string, record *entities.MainUpdate) error {
	return m.handleV2DirectData(remote_addr, record.Basic, record, DIRECT_MAIN_UPDATE_RECORD_COLLECTION_NAME)
}

func (m *MongoDBManager) InsertMainAnnouncementRecord(remote_addr string, record *entities.MainAnnouncement) error {
	return m.handleV2DirectData(remote_addr, record.Basic, record, DIRECT_MAIN_ANNOUNCEMENT_COLLECTION_NAME)
}

func (m *MongoDBManager) InsertUsageQueryRecord(remote_addr string, record *entities.UsageQuery) error {
	return m.handleV2DirectData(remote_addr, record.Basic, record, DIRECT_USAGE_QUERY_COLLECTION_NAME)
}

func (m *MongoDBManager) InsertUsageEventRecord(remote_addr string, record *entities.UsageEvent) error {
	return m.handleV2DirectData(remote_addr, record.Basic, record, DIRECT_USAGE_EVENT_COLLECTION_NAME)
}

func (m *MongoDBManager) InsertUsageFunctionRecord(remote_addr string, record *entities.UsageFunction) error {
	return m.handleV2DirectData(remote_addr, record.Basic, record, DIRECT_USAGE_FUNCTION_COLLECTION_NAME)
}

func (m *MongoDBManager) InsertPluginInstallRecord(remote_addr string, record *entities.PluginInstall) error {
	return m.handleV2DirectData(remote_addr, record.Basic, record, DIRECT_PLUGIN_INSTALL_COLLECTION_NAME)
}

func (m *MongoDBManager) InsertPluginRemoveRecord(remote_addr string, record *entities.PluginRemove) error {
	return m.handleV2DirectData(remote_addr, record.Basic, record, DIRECT_PLUGIN_REMOVE_COLLECTION_NAME)
}

func (m *MongoDBManager) InsertPluginUpdateRecord(remote_addr string, record *entities.PluginUpdate) error {
	return m.handleV2DirectData(remote_addr, record.Basic, record, DIRECT_PLUGIN_UPDATE_COLLECTION_NAME)
}
