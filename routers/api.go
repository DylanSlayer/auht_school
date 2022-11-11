package routers

import (
	"auht_school/api"
	"github.com/gin-gonic/gin"
)

func CollectApi(r *gin.Engine) {
	group := r.Group("/api")
	group.GET("/rainbow", handleRainbowApi)
}

func handleRainbowApi(c *gin.Context) {
	bow, err := api.ApiRainBow()
	if err != nil {
		c.JSON(201, gin.H{
			"code": 201,
			"msg":  "服务器api请求失败，请联系管理员处理",
		})
		return
	}
	c.JSON(200, gin.H{
		"code": 200,
		"data": bow.Result.Content,
	})

}
