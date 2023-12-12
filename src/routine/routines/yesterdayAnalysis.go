package routines

import (
	"context"
	"time"

	"qcg-center/src/database"
	"qcg-center/src/util"

	"go.mongodb.org/mongo-driver/bson"
)

type YesterdayAnalyzeRoutine struct {
	Cfg   *util.Config
	DBMgr *database.MongoDBManager
}

func (r *YesterdayAnalyzeRoutine) Init(cfg *util.Config, db *database.MongoDBManager) (string, error) {
	r.Cfg = cfg
	r.DBMgr = db

	return "@every 12h", nil
}

func (r *YesterdayAnalyzeRoutine) Run() error {
	// 获取今天零时时间戳
	today := time.Now().Truncate(24 * time.Hour)

	yesterday := today.Add(-24 * time.Hour)

	analysis, err := Calc(yesterday, 24*time.Hour, r.DBMgr)

	if err != nil {
		return err
	}

	// 保存到mongo中
	// 检查是否已有相同begin和duration的记录
	// 有则更新，无则插入
	check, err := r.DBMgr.Client.Database("qcg-center-records").Collection("analysis_daily").Find(context.TODO(), map[string]interface{}{
		"begin":    analysis.Begin,
		"duration": analysis.Duration,
	})

	if err != nil {
		return err
	}

	if check.Next(context.Background()) { // 有记录
		// 更新
		_, err = r.DBMgr.Client.Database("qcg-center-records").Collection("analysis_daily").UpdateOne(context.TODO(), map[string]interface{}{
			"begin":    analysis.Begin,
			"duration": analysis.Duration,
		}, bson.M{
			"$set": analysis,
		})
	} else {
		// 插入
		_, err = r.DBMgr.Client.Database("qcg-center-records").Collection("analysis_daily").InsertOne(context.TODO(), analysis)
	}

	return err
}