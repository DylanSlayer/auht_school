package routers

import (
	"auht_school/middleware"
	"auht_school/utils"
	"github.com/gin-gonic/gin"
)

func SetRouters(r *gin.Engine) {

	//跨域中间件
	r.Use(middleware.CORSMiddleware(), middleware.LoginCheckMiddleware())
	utils.InitClient()
	//通用路由
	commonRouters(r)
	//用于登录的router
	LoginRouter(r)
	//课程表相关
	SetCurriculum(r)

	//考试相关
	CollectExaminations(r)
	//第三方api
	CollectApi(r)
}
