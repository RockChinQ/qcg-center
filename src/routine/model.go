package routine

import (
	"qcg-center/src/database"
	"qcg-center/src/util"
)

type IRoutine interface {
	Init(cfg *util.Config, dbmgr *database.MongoDBManager) (string, error) // 初始化，返回cron表达式
	Run() error                                                            // 运行事务
}
