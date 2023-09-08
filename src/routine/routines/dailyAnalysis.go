package routines

import (
	"context"
	"time"

	"qcg-center/src/database"
)

type DailyAnalysis struct {
	Begin             time.Time      `bson:"begin"`
	Duration          int            `bson:"duration"`
	UsageCount        int            `bson:"usage_count"`
	ActiveHostCount   int            `bson:"active_host_count"`
	NewHostCount      int            `bson:"new_host_count"`
	ServiceUsageCount map[string]int `bson:"service_usage_count"`
}

// 计算给定时间段的以下数据：
// - 使用量记录数
// - 活跃主机数
// - 新增主机数(ever)
// - 每种服务的使用量记录数
func Calc(begin time.Time, duration time.Duration, dbmgr *database.MongoDBManager) (*DailyAnalysis, error) {
	result := &DailyAnalysis{
		Begin:    begin,
		Duration: int(duration.Seconds()),
	}

	// 计算使用量记录数
	recnum, err := dbmgr.Client.Database("qcg-center-records").Collection("qchatgpt-usage").CountDocuments(context.TODO(), map[string]interface{}{
		"timestamp": map[string]interface{}{
			"$gte": int64(begin.Unix()),
			"$lt":  int64(begin.Add(duration).Unix()),
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
			"$gte": int64(begin.Unix()),
			"$lt":  int64(begin.Add(duration).Unix()),
		},
	})

	if err != nil {
		return nil, err
	}

	result.ActiveHostCount = len(acthost)

	// fmt.Println("ActiveHostCount:", result.ActiveHostCount)

	// 计算新增主机数
	// 以remote_addr字段去重
	// 判断此时间段每个remote_addr在此之前是否有记录
	// 有则不计入新增主机数
	count := 0
	for _, host := range acthost {
		check, err := dbmgr.Client.Database("qcg-center-records").Collection("qchatgpt-usage").CountDocuments(context.TODO(), map[string]interface{}{
			"remote_addr": host,
			"timestamp": map[string]interface{}{
				"$lt": int64(begin.Unix()),
			},
		})

		if err != nil {
			return nil, err
		}

		if check == 0 {
			count++
		}
	}

	result.NewHostCount = count

	// 计算每种服务的使用量记录数
	// 以service_name进行聚合
	sercount, err := dbmgr.Client.Database("qcg-center-records").Collection("qchatgpt-usage").Aggregate(context.TODO(), []map[string]interface{}{
		{
			"$match": map[string]interface{}{
				"timestamp": map[string]interface{}{
					"$gte": int64(begin.Unix()),
					"$lt":  int64(begin.Add(duration).Unix()),
				},
			},
		},
		{
			"$group": map[string]interface{}{
				"_id": "$service_name",
				"count": map[string]interface{}{
					"$sum": 1,
				},
			},
		},
	})

	if err != nil {
		return nil, err
	}

	result.ServiceUsageCount = make(map[string]int)
	for sercount.Next(context.Background()) {
		var res map[string]interface{}
		err = sercount.Decode(&res)

		if err != nil {
			return nil, err
		}

		result.ServiceUsageCount[res["_id"].(string)] = int(res["count"].(int32))
	}

	return result, nil
}
