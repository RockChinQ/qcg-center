package api

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func GrafanaRoot(m *WebAPI) func(c *gin.Context) {
	return func(c *gin.Context) {
		c.JSON(200, gin.H{
			"code": 0,
			"msg":  "ok",
		})
	}
}

func GrafanaTodayUsageStatic(m *WebAPI) func(c *gin.Context) {
	return func(c *gin.Context) {
		result, err := m.dbmgr.GetTodayUsageStatic()

		if err != nil {
			c.JSON(500, gin.H{
				"code": -1,
				"msg":  err.Error(),
			})
			return
		}

		c.JSON(200, gin.H{
			"code": 0,
			"msg":  "ok",
			"data": result,
		})
	}
}

func GrafanaRecentDaysUsageTrend(m *WebAPI) func(c *gin.Context) {
	return func(c *gin.Context) {
		dayParam := c.Query("day")
		day, err := strconv.Atoi(dayParam)

		if err != nil {
			c.JSON(400, gin.H{
				"code": -1,
				"msg":  err.Error(),
			})
			return
		}

		result, err := m.dbmgr.GetRecentDaysUsageTrend(day)

		if err != nil {
			c.JSON(500, gin.H{
				"code": -1,
				"msg":  err.Error(),
			})
			return
		}

		c.JSON(200, gin.H{
			"code": 0,
			"msg":  "ok",
			"data": result,
		})
	}
}
