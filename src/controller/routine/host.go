package routine

import (
	"log"
	"time"

	"qcg-center/src/database"
	"qcg-center/src/util"

	"github.com/robfig/cron/v3"
)

var Routines map[string]IRoutine
var cronHost *cron.Cron

// 注册routine
func Register(name string, routine IRoutine) {
	if Routines == nil {
		Routines = make(map[string]IRoutine)
	}
	Routines[name] = routine
}

func ScheduleAll(cfg *util.Config, dbmgr *database.MongoDBManager) {
	if cronHost == nil {
		loc, err := time.LoadLocation("Asia/Shanghai")
		if err != nil {
			panic(err)
		}

		cronHost = cron.New(cron.WithLocation(loc))
	}

	for name, routine := range Routines {

		// 克隆name和routine
		cname := name
		croutine := routine

		log.Println("Scheduling routine", cname)
		cronExpr, err := routine.Init(cfg, dbmgr)
		if err != nil {
			log.Printf("routine %s init failed: %s", cname, err)
			panic(err)
		}

		cronHost.AddFunc(cronExpr, func() {
			log.Println("Running routine", cname)
			err := croutine.Run()
			if err != nil {
				log.Printf("routine %s run failed: %s", cname, err)
			} else {
				log.Println("Routine", cname, "finished")
			}
		})
	}
	cronHost.Start()
}
