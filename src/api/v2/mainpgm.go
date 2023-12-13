package v2

import (
	"qcg-center/src/database"
	"qcg-center/src/entities"

	"github.com/gin-gonic/gin"
)

func MainUpdate(db *database.IDatabaseManager) func(c *gin.Context) {
	return func(c *gin.Context) {
		handleRequest[entities.MainUpdate](c, db, (*db).InsertMainUpdateRecord)
	}
}

func MainAnnouncement(db *database.IDatabaseManager) func(c *gin.Context) {
	return func(c *gin.Context) {
		handleRequest[entities.MainAnnouncement](c, db, (*db).InsertMainAnnouncementRecord)
	}
}
