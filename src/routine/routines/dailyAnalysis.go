package routines

import (
	"context"
	"log"
	"time"

	"qcg-center/src/database"
)

type DailyAnalysis struct {
	Begin           time.Time `bson:"begin"`
	Duration        int       `bson:"duration"`
	UsageCount      int       `bson:"usage_count"`
	ActiveHostCount int       `bson:"active_host_count"`
	NewHostCount    int       `bson:"new_host_count"`
}

// 计算给定时间段的以下数据：
// - 使用量记录数
// - 活跃主机数
// - 新增主机数(ever)
func Calc(begin time.Time, duration time.Duration, dbmgr *database.MongoDBManager) (*DailyAnalysis, error) {

	result := &DailyAnalysis{
		Begin:    begin,
		Duration: int(duration.Seconds()),
	}

	// 计算使用量记录数
	recnum, err := dbmgr.Client.Database("qcg-center-records").Collection("qchatgpt-usage").CountDocuments(context.TODO(), map[string]interface{}{
		"timestamp": map[string]interface{}{
			"$gte": int64(begin.Add(-8 * time.Hour).Unix()),
			"$lt":  int64(begin.Add(-8 * time.Hour).Add(duration).Unix()),
		},
	})

	if err != nil {
		return nil, err
	}

	result.UsageCount = int(recnum)

	// fmt.Println("UsageCount:", result.UsageCount)

	// 计算活跃主机数
	// 以remote_addr字段去重
	acthost, err := dbmgr.Client.Database("qcg-center-records").Collection("qchatgpt-usage").Distinct(context.TODO(), "remote_addr", map[string]interface{}{
		"timestamp": map[string]interface{}{
			"$gte": int64(begin.Add(-8 * time.Hour).Unix()),
			"$lt":  int64(begin.Add(-8 * time.Hour).Add(duration).Unix()),
		},
	})

	if err != nil {
		return nil, err
	}

	result.ActiveHostCount = len(acthost)

	// fmt.Println("ActiveHostCount:", result.ActiveHostCount)

	// 计算新增主机数
	// 从analysis_usage_remote_addrs查找created_at在此时段的记录数量

	newcount, err := dbmgr.Client.Database("qcg-center-records").Collection("analysis_usage_remote_addrs").CountDocuments(context.TODO(), map[string]interface{}{
		"created_at": map[string]interface{}{
			"$gte": begin.Add(-8 * time.Hour),
			"$lt":  begin.Add(-8 * time.Hour).Add(duration),
		},
	})

	if err != nil {
		return nil, err
	}

	result.NewHostCount = int(newcount)

	// 输出格式化后的结果
	// 包含：开始时间、时长、使用量记录数、活跃主机数、新增主机数
	log.Printf("DailyAnalysis: %s, %d, %d, %d, %d\n", result.Begin.Format("2006-01-02 15:04:05"), result.Duration, result.UsageCount, result.ActiveHostCount, result.NewHostCount)

	return result, nil
}
