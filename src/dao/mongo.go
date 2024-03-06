package dao

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"qcg-center/src/entities/dao"
	"qcg-center/src/entities/dto"
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
	Client   *mongo.Client
	Database string
}

func NewMongoDBManager(Uri string, Database string) (*MongoDBManager, error) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(Uri))
	if err != nil {
		return nil, err
	}
	return &MongoDBManager{Client: client, Database: Database}, nil
}

func (m *MongoDBManager) InsertIPGeoInfo(record dao.IPGeoInfoDAO) error {

	coll_name := ANALYSIS_IP_COLLECTION_NAME

	coll := m.Client.Database(m.Database).Collection(coll_name)

	// 检查是否存在相同的IP记录
	check, err := coll.Find(context.TODO(), bson.M{"ip": record.IP})

	if err != nil {
		return err
	}

	if check.Next(context.TODO()) {
		return nil
	} else {
		_, err := coll.InsertOne(context.TODO(), record)
		return err
	}

}

func (m *MongoDBManager) InsertIdentifierTuple(record dao.IdentifierTupleDAO) error {

	coll_name := ANALYSIS_HOST_INSTANCE_IP_COLLECTION_NAME

	coll := m.Client.Database(m.Database).Collection(coll_name)

	// 检查是否存在相同的记录
	check, err := coll.Find(context.TODO(), bson.M{"host_id": record.HostID, "instance_id": record.InstanceID, "ip": record.IP})

	if err != nil {
		return err
	}

	if check.Next(context.TODO()) {
		return nil
	} else {
		_, err := coll.InsertOne(context.TODO(), record)
		return err
	}

}

func (m *MongoDBManager) InsertRecord(record dao.CommonRecordDAO) error {
	// 根据 record.data 位于 DTO 包中的类型确定插入的集合
	// 由于 record.data 是 interface{} 类型，需要使用反射获取其类型
	coll_name := ""

	switch record.Data.(type) {
	case dto.MainUpdateDTO:
		coll_name = DIRECT_MAIN_UPDATE_RECORD_COLLECTION_NAME
	case dto.MainAnnouncementDTO:
		coll_name = DIRECT_MAIN_ANNOUNCEMENT_COLLECTION_NAME
	case dto.UsageQueryDTO:
		coll_name = DIRECT_USAGE_QUERY_COLLECTION_NAME
	case dto.UsageEventDTO:
		coll_name = DIRECT_USAGE_EVENT_COLLECTION_NAME
	case dto.UsageFunctionDTO:
		coll_name = DIRECT_USAGE_FUNCTION_COLLECTION_NAME
	case dto.PluginInstallDTO:
		coll_name = DIRECT_PLUGIN_INSTALL_COLLECTION_NAME
	case dto.PluginRemoveDTO:
		coll_name = DIRECT_PLUGIN_REMOVE_COLLECTION_NAME
	case dto.PluginUpdateDTO:
		coll_name = DIRECT_PLUGIN_UPDATE_COLLECTION_NAME
	default:
		return errors.New("unknown record type")
	}

	coll := m.Client.Database(m.Database).Collection(coll_name)

	_, err := coll.InsertOne(context.TODO(), record)

	return err
}

func (m *MongoDBManager) CountUniqueValueInDuration(coll_name string, field_name string, start_time time.Time, end_time time.Time, time_field_name string) (int, error) {
	coll := m.Client.Database(m.Database).Collection(coll_name)

	pipeline := bson.A{
		bson.D{{Key: "$match", Value: bson.D{{Key: time_field_name, Value: bson.D{{Key: "$gte", Value: start_time}, {Key: "$lte", Value: end_time}}}}}},
		bson.D{{Key: "$group", Value: bson.D{{Key: "_id", Value: "$" + field_name}}}},
		bson.D{{Key: "$group", Value: bson.D{{Key: "_id", Value: nil}, {Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}}}}},
	}

	cursor, err := coll.Aggregate(context.TODO(), pipeline)

	if err != nil {
		return 0, err
	}

	var result []bson.M

	if err = cursor.All(context.Background(), &result); err != nil {
		return 0, err
	}

	if len(result) == 0 {
		return 0, nil
	}

	return int(result[0]["count"].(int32)), nil
}
