// 实时数据
package view

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"qcg-center/src/dao"
	"qcg-center/src/service/view"
	"qcg-center/src/util"

	"time"
)

// 返回通用的唯一值计数处理器函数
func CommonUniqueValueCountingGeneric(
	sv *view.RealTimeDataService,
	coll_name string,
	field_name string,
	time_field_name string,
) func(c *gin.Context) {
	return func(c *gin.Context) {
		minute_param := c.Query("minute")

		start_time_param := c.Query("start_time")
		end_time_param := c.Query("end_time")

		period := c.Query("period")
		period_amount := c.Query("period_amount")
		period_offset := c.Query("period_offset")

		timezone := c.Query("timezone")

		// if minute_param == "" && (start_time_param == "" || end_time_param == "") {
		// 	c.JSON(400, gin.H{"error": "missing parameter"})
		// 	return
		// }

		var start_time time.Time
		var end_time time.Time

		if minute_param != "" {
			minute, err := strconv.Atoi(minute_param)
			if err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}

			now := util.GetCSTTime()

			start_time = now.Add(-1 * time.Duration(minute) * time.Minute)
			end_time = now
		} else if start_time_param != "" && end_time_param != "" {
			var err error
			layout := "2006-01-02 15:04:05"

			location := time.Local

			if timezone != "" {
				var err error
				location, err = time.LoadLocation(timezone)
				if err != nil {
					c.JSON(400, gin.H{"error": err.Error()})
					return
				}
			}

			start_time, err = time.ParseInLocation(layout, start_time_param, location)
			if err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}

			end_time, err = time.ParseInLocation(layout, end_time_param, location)
			if err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}

			if start_time.After(end_time) {
				c.JSON(400, gin.H{"error": "start time should be before end time"})
				return
			}
		} else if period != "" && period_amount != "" && period_offset != "" {
			// 取整
			period_amount_int, err := strconv.Atoi(period_amount)
			if err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}

			period_offset_int, err := strconv.Atoi(period_offset)
			if err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}

			// 获取开始时间的整period时间点
			start_time, end_time = util.GetCSTFixedPeriodTime(period, period_amount_int, period_offset_int)

		} else {
			c.JSON(400, gin.H{"error": "missing parameter"})
			return
		}

		count, err := sv.CountUniqueValueInDuration(coll_name, field_name, start_time, end_time, time_field_name)

		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{
			"code": 0,
			"msg":  "ok",
			"data": gin.H{
				"count": count,
			},
		})
	}
}

func BindPath(r *gin.Engine, sv *view.RealTimeDataService) {
	// unique value counting
	r.GET("/api/v2/view/realtime/count_query", CommonUniqueValueCountingGeneric(sv, dao.DIRECT_USAGE_QUERY_COLLECTION_NAME, "data.basic.rid", "time"))
	r.GET("/api/v2/view/realtime/count_event", CommonUniqueValueCountingGeneric(sv, dao.DIRECT_USAGE_EVENT_COLLECTION_NAME, "data.basic.rid", "time"))
	r.GET("/api/v2/view/realtime/count_function", CommonUniqueValueCountingGeneric(sv, dao.DIRECT_USAGE_FUNCTION_COLLECTION_NAME, "data.basic.rid", "time"))
	r.GET("/api/v2/view/realtime/count_plugin_install", CommonUniqueValueCountingGeneric(sv, dao.DIRECT_PLUGIN_INSTALL_COLLECTION_NAME, "data.basic.rid", "time"))

	r.GET("/api/v2/view/realtime/count_active_instance", CommonUniqueValueCountingGeneric(sv, dao.DIRECT_USAGE_QUERY_COLLECTION_NAME, "data.basic.instance_id", "time"))
	r.GET("/api/v2/view/realtime/count_active_host", CommonUniqueValueCountingGeneric(sv, dao.DIRECT_USAGE_QUERY_COLLECTION_NAME, "data.basic.host_id", "time"))
}
