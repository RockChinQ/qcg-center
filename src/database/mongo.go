package database

import (
	"context"
	"time"

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

func (m *MongoDBManager) GetTodayUsageStatic() (TodayUsageStatic, error) {
	var todayUsageStatic TodayUsageStatic

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
