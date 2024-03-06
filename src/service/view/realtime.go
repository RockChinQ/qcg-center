// 实时数据展示
package view

import (
	"qcg-center/src/dao"
	"time"
)

type RealTimeDataService struct {
	db dao.IDatabaseManager
}

func NewRealTimeDataService(db dao.IDatabaseManager) *RealTimeDataService {
	return &RealTimeDataService{
		db: db,
	}
}

// 获取一段时间内某字段的唯一值的数量
func (s *RealTimeDataService) CountUniqueValueInDuration(
	coll_name string,
	field_name string,
	start_time time.Time,
	end_time time.Time,
	time_field_name string,
) (int, error) {
	// TODO
	return s.db.CountUniqueValueInDuration(coll_name, field_name, start_time, end_time, time_field_name)
}
