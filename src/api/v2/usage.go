package v2

import (
	"qcg-center/src/database"
	"qcg-center/src/entities"

	"github.com/gin-gonic/gin"
)

func UsageQuery(db *database.IDatabaseManager) func(c *gin.Context) {
	return func(c *gin.Context) {
		handleRequest[entities.UsageQuery](c, db, (*db).InsertUsageQueryRecord)
	}
}

func UsageEvent(db *database.IDatabaseManager) func(c *gin.Context) {
	return func(c *gin.Context) {
		handleRequest[entities.UsageEvent](c, db, (*db).InsertUsageEventRecord)
	}
}

func UsageFunction(db *database.IDatabaseManager) func(c *gin.Context) {
	return func(c *gin.Context) {
		handleRequest[entities.UsageFunction](c, db, (*db).InsertUsageFunctionRecord)
	}
}
