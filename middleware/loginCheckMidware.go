package middleware

import (
	"auht_school/utils"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

func LoginCheckMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.FullPath() == "/login" {
			//登录时不进行拦截
			c.Next()
			return
		}
		targetUrl := "https://vpncas.ahut.edu.cn/http/77726476706e69737468656265737421fae05988693160456a468ca88d1b203b/jsxsd/framework/xsMain.jsp"
		request, err := utils.GetRequest(c, "GET", targetUrl, "", nil)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, "服务器内部错误，无法创建请求，请稍后重试或联系管理员处理")
			return
		}
		response, err := utils.Client.Do(request)
		cookies := response.Request.Cookies()
		log.Println("cookies:-------------", cookies)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, "服务器无法请求到目标服务器，可能是学校服务器出现故障，或服务器内部出现异常，请稍后重试或联系管理员处理")
			return
		}
		all, _ := ioutil.ReadAll(response.Body)
		defer response.Body.Close()
		body := string(all)
		title := regexp.MustCompile(`<title>(.*?)</title>`)
		findString := title.FindString(body)
		log.Println("logString   ", findString)
		if findString == "<title>Login - 安徽工业大学</title>" {
			//用户没有登录或登录失效
			c.AbortWithStatusJSON(200, gin.H{
				"code": 211,
				"msg":  "用户登录消息过期请重新登录",
			})
			return
		}
		c.Next()
	}
}
