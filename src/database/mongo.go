package database

import (
	"context"
	"time"

	"qcg-center/src/entities"
	"qcg-center/src/util"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDBManager struct {
	Cfg *util.Config

	Client *mongo.Client
}

func (m *MongoDBManager) Connect() error {
	uri := m.Cfg.Database.Params["uri"]

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
			"created_at": time.Unix(report.Timestamp, 0),
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
			"created_at":  time.Unix(usage.Timestamp, 0),
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
	today = time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, time.UTC)

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
	today = time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, time.UTC)

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

// v2
func (m *MongoDBManager) InsertMainUpdateRecord(remote_addr string, record *entities.MainUpdate) error {
	return nil
}

func (m *MongoDBManager) InsertMainAnnouncementRecord(remote_addr string, record *entities.MainAnnouncement) error {
	return nil
}

func (m *MongoDBManager) InsertUsageQueryRecord(remote_addr string, record *entities.UsageQuery) error {

	return nil
}

func (m *MongoDBManager) InsertUsageEventRecord(remote_addr string, record *entities.UsageEvent) error {

	return nil
}

func (m *MongoDBManager) InsertUsageFunctionRecord(remote_addr string, record *entities.UsageFunction) error {

	return nil
}

func (m *MongoDBManager) InsertPluginInstallRecord(remote_addr string, record *entities.PluginInstall) error {

	return nil
}

func (m *MongoDBManager) InsertPluginRemoveRecord(remote_addr string, record *entities.PluginRemove) error {

	return nil
}

func (m *MongoDBManager) InsertPluginUpdateRecord(remote_addr string, record *entities.PluginUpdate) error {

	return nil
}
