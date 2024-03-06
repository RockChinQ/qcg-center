// 实时数据
package view

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"

	"qcg-center/src/dao"
	"qcg-center/src/service/view"
	"qcg-center/src/util"

	"time"
)

func GetTimeRange(c *gin.Context) (time.Time, time.Time, error) {

	minute_param := c.Query("minute")

	start_time_param := c.Query("start_time")
	end_time_param := c.Query("end_time")

	period := c.Query("period")
	period_amount := c.Query("period_amount")
	period_offset := c.Query("period_offset")

	timezone := c.Query("timezone")

	var start_time time.Time
	var end_time time.Time

	if minute_param != "" {
		minute, err := strconv.Atoi(minute_param)
		if err != nil {
			return time.Time{}, time.Time{}, err
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
				return time.Time{}, time.Time{}, err
			}
		}

		start_time, err = time.ParseInLocation(layout, start_time_param, location)
		if err != nil {
			return time.Time{}, time.Time{}, err
		}

		end_time, err = time.ParseInLocation(layout, end_time_param, location)
		if err != nil {
			return time.Time{}, time.Time{}, err
		}

		if start_time.After(end_time) {
			return time.Time{}, time.Time{}, err
		}
	} else if period != "" && period_amount != "" && period_offset != "" {
		// 取整
		period_amount_int, err := strconv.Atoi(period_amount)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return time.Time{}, time.Time{}, err
		}

		period_offset_int, err := strconv.Atoi(period_offset)
		if err != nil {
			return time.Time{}, time.Time{}, err
		}

		// 获取开始时间的整period时间点
		start_time, end_time = util.GetCSTFixedPeriodTime(period, period_amount_int, period_offset_int)

	} else {
		return time.Time{}, time.Time{}, errors.New("invalid time range")
	}

	return start_time, end_time, nil
}

// 返回通用的唯一值计数处理器函数
func CommonUniqueValueCountingGeneric(
	sv *view.RealTimeDataService,
	coll_name string,
	field_name string,
	time_field_name string,
) func(c *gin.Context) {
	return func(c *gin.Context) {
		start_time, end_time, err := GetTimeRange(c)

		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
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

func CommonAggragationValueAmountGeneric(
	sv *view.RealTimeDataService,
	coll_name string,
	field_name string,
	time_field_name string,
) func(c *gin.Context) {
	return func(c *gin.Context) {
		start_time, end_time, err := GetTimeRange(c)

		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		result, err := sv.AggregationValueAmountInDuration(coll_name, field_name, start_time, end_time, time_field_name)

		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{
			"code": 0,
			"msg":  "ok",
			"data": result,
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

	r.GET("/api/v2/view/realtime/get_os_amount", CommonAggragationValueAmountGeneric(sv, dao.DIRECT_USAGE_QUERY_COLLECTION_NAME, "data.basic.platform", "time"))
	r.GET("/api/v2/view/realtime/get_version_amount", CommonAggragationValueAmountGeneric(sv, dao.DIRECT_USAGE_QUERY_COLLECTION_NAME, "data.basic.semantic_version", "time"))
	r.GET("/api/v2/view/realtime/get_platform_amount", CommonAggragationValueAmountGeneric(sv, dao.DIRECT_USAGE_QUERY_COLLECTION_NAME, "data.runtime.msg_source", "time"))
	// data.session_info.type
	r.GET("/api/v2/view/realtime/get_session_type_amount", CommonAggragationValueAmountGeneric(sv, dao.DIRECT_USAGE_QUERY_COLLECTION_NAME, "data.session_info.type", "time"))
	// data.query_info.model_name
	r.GET("/api/v2/view/realtime/get_model_name_amount", CommonAggragationValueAmountGeneric(sv, dao.DIRECT_USAGE_QUERY_COLLECTION_NAME, "data.query_info.model_name", "time"))
}
