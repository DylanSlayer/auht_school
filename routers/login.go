package routers

import (
	"auht_school/model"
	"auht_school/utils"
	"bytes"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

// LoginRouter 用户登录,获取登录后的cookie
func LoginRouter(r *gin.Engine) {
	r.POST("/login", loginHandler)
	r.GET("/getUserInfo", getUserInfo)
}

//登录检查，查看登录信息是否有效,没问题返回学生基本信息
func getUserInfo(c *gin.Context) {

	targetUrl := "http://jwxt.ahut.edu.cn/jsxsd/framework/xsMain_new.jsp?t1=1"
	request, err := utils.GetRequest(c, "GET", targetUrl, "", nil)
	if err != nil {
		c.String(http.StatusInternalServerError, "服务器内部错误，无法创建请求，请稍后重试或联系管理员处理")
		return
	}
	response, err := utils.Client.Do(request)
	if err != nil {
		c.String(http.StatusBadRequest, "服务器无法请求到目标服务器，可能是学校服务器出现故障，或服务器内部出现异常，请稍后重试或联系管理员处理")
		return
	}
	all, _ := ioutil.ReadAll(response.Body)
	html := string(all)

	html = strings.Replace(html, "\n", "", -1)
	html = strings.Replace(html, "\r", "", -1)
	html = strings.Replace(html, "\t", "", -1)
	//获取用户信息
	userInfoCompile := regexp.MustCompile(`<div class="middletopdwxxcont">(.*?)</div>`)
	userInfo := userInfoCompile.FindAllString(html, -1)
	if len(userInfo) < 6 {
		c.JSON(200, gin.H{
			"code": 201,
			"data": model.UserInfo{Name: "服务器异常，请稍后重试"},
			"msg":  "无法获取用户数据",
		})
	}
	data := model.UserInfo{
		Name:       userInfo[1][31 : len(userInfo[1])-6],
		Id:         userInfo[2][31 : len(userInfo[2])-6],
		Academy:    userInfo[3][31 : len(userInfo[3])-6],
		Profession: userInfo[4][31 : len(userInfo[4])-6],
		Class:      userInfo[5][31 : len(userInfo[5])-6],
	}
	log.Println(data)

	//log.Println(findString)
	c.JSON(200, gin.H{
		"code": 200,
		"data": data,
	})
}

//登录并返回cookie
func loginHandler(c *gin.Context) {
	//获取userAccount和encoded
	loginForm := model.LoginForm{}
	err := c.ShouldBind(&loginForm)
	if err != nil {
		c.String(http.StatusBadRequest, "登录失败,用户信息绑定错误")
		return
	}
	//像服务器发送登录请求
	targetUrl := "http://jwxt.ahut.edu.cn/jsxsd/xk/LoginToXk"
	log.Println("loginForm:====> {}", loginForm)

	//创建form数据
	// 用url.values方式构造form-data参数
	formValues := url.Values{}
	formValues.Set("userAccount", loginForm.UserAccount)
	formValues.Set("encoded", loginForm.Encoded)
	formDataStr := formValues.Encode()
	formDataBytes := []byte(formDataStr)
	formBytesReader := bytes.NewReader(formDataBytes)
	//创建请求
	req, err := utils.GetRequest(c, "POST", targetUrl, "application/x-www-form-urlencoded", formBytesReader)
	if err != nil {
		c.JSON(200, gin.H{
			"code": 433,
			"msg":  "服务器故障，无法创建请求，请联系管理员进行处理",
		})
		return
	}
	//发送请求
	resp, err := utils.Client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		c.String(200, "登录失败，服务器发送的请求失败")
		return
	}

	all, _ := ioutil.ReadAll(resp.Body)
	body := string(all)
	//判断是否登录成功
	loginTips_compile := regexp.MustCompile(`<span style="color: red; font-size:16px">(.*?)</span>`)
	loginTips := loginTips_compile.FindString(body)
	log.Println("登录html:", body)
	log.Println("登录消息：", loginTips)
	if len(loginTips) > len("<span style=\"color: red; font-size:16px\"></span>") {
		//判定为登录失败返回登录失败的消息
		c.JSON(200, gin.H{
			"code": 433,
			"msg":  loginTips[41 : len(loginTips)-7],
		})
		return
	}
	cookie := resp.Cookies()
	//将cookie返回
	c.JSON(200, model.Result{
		Code: http.StatusOK,
		Data: cookie,
		Msg:  "登录成功",
	})
}
