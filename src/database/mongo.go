package database

import (
	"context"

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

	return err
}

func (m *MongoDBManager) StoreQChatGPTUsage(usage *QChatGPTUsage) error {
	coll := m.Client.Database("qcg-center-records").Collection("qchatgpt-usage")

	_, err := coll.InsertOne(context.TODO(), usage)

	return err
}
