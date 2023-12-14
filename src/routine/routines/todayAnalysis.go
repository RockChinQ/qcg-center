package routines

import (
	"context"
	"time"

	"qcg-center/src/database"
	"qcg-center/src/util"

	"go.mongodb.org/mongo-driver/bson"
)

type TodayAnalyzeRoutine struct {
	Cfg   *util.Config
	DBMgr *database.MongoDBManager
}

func (r *TodayAnalyzeRoutine) Init(cfg *util.Config, db *database.MongoDBManager) (string, error) {
	r.Cfg = cfg
	r.DBMgr = db

	return "@every 15m", nil
	// return "45 19 * * *", nil
}

func (r *TodayAnalyzeRoutine) Run() error {
	// 统一 UTC
	today := time.Now().UTC()

	// 今天的0点
	today = time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, util.GetCSTTimeLocation())

	analysis, err := Calc(today, 24*time.Hour, r.DBMgr)

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
