package api

import (
	database "qcg-center/src/database"
)

type IAPIManager interface {

	// 初始化API管理器
	Init(dbmgr database.IDatabaseManager) error

	// 启动API服务
	Serve() error
}
