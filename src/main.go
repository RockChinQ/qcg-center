package main

import (
	"log"
	"os"
	"qcg-center/src/controller/api"
	"qcg-center/src/controller/routine"
	"qcg-center/src/controller/routine/routines"
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

	dbmgr := &database.MongoDBManager{
		Cfg: cfg,
	}

	err = dbmgr.Connect()

	if err != nil {
		log.Println("Failed to connect to database")
		panic(err)
	}

	apimgr := &api.WebAPI{
		Cfg: cfg,
	}

	err = apimgr.Init(dbmgr)

	if err != nil {
		log.Println("Failed to initialize API manager")
		panic(err)
	}

	// 初始化routines
	InitializeRoutines(cfg, dbmgr)
	log.Printf("Routines scheduled.")

	err = apimgr.Serve()

	if err != nil {
		log.Println("Failed to serve API")
		panic(err)
	}
}

func InitializeRoutines(cfg *util.Config, db *database.MongoDBManager) {
	// 注册routines
	routine.Register("TodayAnalysis", &routines.TodayAnalyzeRoutine{})
	routine.Register("YesterdayAnalysis", &routines.YesterdayAnalyzeRoutine{})

	// 启动routines
	routine.ScheduleAll(cfg, db)
}
