package main

import (
	"log"
	"os"
	"qcg-center/src/controller/api"
	"qcg-center/src/dao"
	serviceRecord "qcg-center/src/service/record"
	"qcg-center/src/util"
)

func main() {
	log.Println("Launching QChatGPT Center...")

	// 检查并加载配置文件
	_, created, err := util.EnsureConfigFile()

	if err != nil {
		log.Println("Failed to load config.yaml")
		panic(err)
	}

	if created {
		log.Println("Created config.yaml, please edit it and restart the program")
		os.Exit(0)
	}

	cfg, err := util.LoadConfig()

	if err != nil {
		log.Println("Failed to load config.yaml")
		panic(err)
	}

	// 初始化数据库管理器
	dbmgr, err := dao.NewMongoDBManager(cfg.Database.Params["uri"], cfg.Database.Params["database"])

	if err != nil {
		log.Println("Failed to connect to database")
		panic(err)
	}

	// 初始化服务
	svRecord := serviceRecord.NewRecordService(dbmgr)

	// 初始化API管理器
	apimgr := api.NewWebAPI(svRecord, cfg.API.Port, cfg.API.Listen)

	// 初始化routines
	// InitializeRoutines(cfg, dbmgr)
	// log.Printf("Routines scheduled.")

	err = apimgr.Serve()

	if err != nil {
		log.Println("Failed to serve API")
		panic(err)
	}
}
