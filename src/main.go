package main

import (
	"log"
	"os"
	"qcg-center/src/api"
	"qcg-center/src/database"
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
	var dbmgr database.IDatabaseManager

	dbmgr = &database.MongoDBManager{
		Cfg: cfg,
	}

	err = dbmgr.Connect()

	if err != nil {
		log.Println("Failed to connect to database")
		panic(err)
	}

	var apimgr api.IAPIManager

	apimgr = &api.WebAPI{
		Cfg: cfg,
	}

	err = apimgr.Init(dbmgr)

	if err != nil {
		log.Println("Failed to initialize API manager")
		panic(err)
	}

	err = apimgr.Serve()

	if err != nil {
		log.Println("Failed to serve API")
		panic(err)
	}
}
