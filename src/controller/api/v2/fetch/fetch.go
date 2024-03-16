package fetch

import (
	"github.com/gin-gonic/gin"

	fetchsv "qcg-center/src/service/fetch"
)

func BindPath(r *gin.Engine, sv *fetchsv.FetchService) {
	r.GET("/api/v2/fetch/model_list", func(ctx *gin.Context) {
		version := ctx.Query("version")

		modelList, err := sv.FetchModelList(ctx, version)

		if err != nil {
			ctx.JSON(500, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(200, gin.H{
			"code": 0,
			"msg":  "ok",
			"data": gin.H{
				"list": modelList,
			},
		})
	})
}
