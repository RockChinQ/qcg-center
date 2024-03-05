package v2

import (
	"qcg-center/src/database"
	"qcg-center/src/entities"

	"github.com/gin-gonic/gin"
)

func PluginInstall(db *database.IDatabaseManager) func(c *gin.Context) {
	return func(c *gin.Context) {
		handleRequest[entities.PluginInstall](c, db, (*db).InsertPluginInstallRecord)
	}
}

func PluginRemove(db *database.IDatabaseManager) func(c *gin.Context) {
	return func(c *gin.Context) {
		handleRequest[entities.PluginRemove](c, db, (*db).InsertPluginRemoveRecord)
	}
}

func PluginUpdate(db *database.IDatabaseManager) func(c *gin.Context) {
	return func(c *gin.Context) {
		handleRequest[entities.PluginUpdate](c, db, (*db).InsertPluginUpdateRecord)
	}
}
