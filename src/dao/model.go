package dao

import (
	"qcg-center/src/entities/dao"
	"qcg-center/src/entities/dto"
	"time"
)

type IDatabaseManager interface {

	// 插入IP地理信息
	// 需要由实现检查是否已经存在
	InsertIPGeoInfo(record dao.IPGeoInfoDAO) error

	// 插入标识符元组
	// 需要由实现检查是否已经存在
	InsertIdentifierTuple(record dao.IdentifierTupleDAO) error

	// 插入记录数据
	InsertRecord(record dao.CommonRecordDAO) error

	// 统计一段时间内某字段的唯一值的数量
	CountUniqueValueInDuration(coll_name string, field_name string, start_time time.Time, end_time time.Time, time_field_name string) (int, error)

	// 聚合一段时间内某个字段的相同值的数量
	AggregationValueAmountInDuration(coll_name string, field_name string, start_time time.Time, end_time time.Time, time_field_name string) (dto.AggregationValueAmountDTO, error)
}
