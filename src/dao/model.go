package dao

import (
	"qcg-center/src/entities/dao"
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
}
