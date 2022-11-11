package routers

import (
	"auht_school/common"
	"auht_school/model"
	"auht_school/utils"
	"bytes"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/url"
	"regexp"
	"strings"
)

func CollectExaminations(r *gin.Engine) {

	r.GET("/getExamInClass", getExamInClass)
	r.GET("/examCommon", getExamCommon)
}

//随堂考试查询
func getExamInClass(c *gin.Context) {
	targetUtl := common.Settings.ServInfo.BaseUrl + "/jsxsd/xsks/xsstk_list"
	//先请求随堂考试获取考试信息
	formValues := url.Values{}
	formValues.Set("xnxqid", utils.GetCurrentTerm())
	formDataStr := formValues.Encode()
	formDataBytes := []byte(formDataStr)
	formBytesReader := bytes.NewReader(formDataBytes)

	request, err := utils.GetRequest(c, "POST", targetUtl, "application/x-www-form-urlencoded", formBytesReader)
	if err != nil {
		c.JSON(200, gin.H{
			"code": 201,
			"msg":  "服务器无法创建请求，请联系管理员或稍后重试",
		})
		return
	}
	response, err := utils.Client.Do(request)
	if err != nil {
		c.JSON(200, gin.H{
			"code": 201,
			"msg":  "无法访问到学校的访问器",
		})
		return
	}
	all, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	if err != nil {
		c.JSON(200, gin.H{
			"code": 201,
			"msg":  "服务器内部错误，无法读取响应体",
		})
	}
	html := string(all)
	dealer := tdTrDealer(html, 10)
	if dealer == nil {
		c.JSON(200, gin.H{
			"code": 200,
			"data": "",
			"msg":  "未查找到相关信息",
		})
		return
	}
	var data []model.ExamInfo
	for i := 0; i < len(dealer); i++ {
		data = append(data, model.ExamInfo{
			Term:      dealer[i][0],
			Code:      dealer[i][1],
			Name:      dealer[i][2],
			Week:      dealer[i][3],
			WeekDay:   dealer[i][4],
			Class:     dealer[i][5],
			Teacher:   dealer[i][6],
			ClassRoom: dealer[i][7],
			Time:      dealer[i][8],
			Type:      dealer[i][9],
		})
	}

	log.Println(data)
	c.JSON(200, gin.H{
		"code": 200,
		"data": data,
	})
}

func getExamCommon(c *gin.Context) {
	targetUrl := common.Settings.ServInfo.BaseUrl + "/jsxsd/xsks/xsksap_list"
	method := "POST"
	contentType := "application/x-www-form-urlencoded"
	formValues := url.Values{}
	formValues.Set("xnxqid", utils.GetCurrentTerm())
	formDataStr := formValues.Encode()
	formDataBytes := []byte(formDataStr)
	formBytesReader := bytes.NewReader(formDataBytes)
	request, err := utils.GetRequest(c, method, targetUrl, contentType, formBytesReader)
	if err != nil {
		c.JSON(200, gin.H{
			"code": 201,
			"msg":  "服务器无法创建请求，请联系管理员或稍后重试",
		})
		return
	}
	do, err := utils.Client.Do(request)
	if err != nil {
		c.JSON(200, gin.H{
			"code": 201,
			"msg":  "无法访问到学校的访问器",
		})
	}
	all, _ := ioutil.ReadAll(do.Body)
	html := string(all)
	dealer := commonExamInfoDealer(html)
	if dealer == nil {
		c.JSON(200, gin.H{
			"code": 200,
			"data": "",
			"msg":  "未查找到相关信息",
		})
		return
	}
	var data []model.CommonExamInfo
	for i := 0; i < len(dealer); i++ {
		data = append(data, model.CommonExamInfo{
			Id:         dealer[i][0],
			Campus:     dealer[i][1],
			ExamCampus: dealer[i][2],
			Session:    dealer[i][3],
			Number:     dealer[i][4],
			Name:       dealer[i][5],
			Teacher:    dealer[i][6],
			Time:       dealer[i][7],
			Class:      dealer[i][8],
			Seat:       dealer[i][9],
		})
	}
	c.JSON(200, gin.H{
		"code": 200,
		"data": data,
	})
}

func commonExamInfoDealer(html string) [][]string {
	html = strings.Replace(html, "\r", "", -1)
	html = strings.Replace(html, "\n", "", -1)
	html = strings.Replace(html, "\t", "", -1)
	// 最后一个是操作中的备注，不需要
	tdCompile1 := regexp.MustCompile(`<td>(.*?)</td>`)
	// 校区，考场校区，考试场次，课程编号，课程名称
	tdCompile2 := regexp.MustCompile(`<td align="left">(.*?)</td>`)

	//查找考试条目个数
	trCompile := regexp.MustCompile(`<tr>(.*?)</tr>`)
	row := len(trCompile.FindAllString(html, -1)) - 1

	if row == 0 {
		return nil
	}
	//序号，授课教师，时间，考场，座位号， 。。。 ， 。。。 ， 。。。
	td1s := tdCompile1.FindAllString(html, -1)
	//校区，考场校区，考试场次，课程编号，课程名称
	td2s := tdCompile2.FindAllString(html, -1)

	//封装有用信息
	target := make([][]string, row)
	for i := range target {
		target[i] = make([]string, 10)
	}
	for i := 0; i < row; i++ {
		target[i][0] = td1s[i*8+0]
		for j := 0; j < 5; j++ {
			target[i][j+1] = td2s[i*5+j]
		}
		for j := 0; j < 4; j++ {
			target[i][j+6] = td1s[i*8+j+1]
		}
	}
	return target
}

//html为待处理文本，col每一行有多少数据
func tdTrDealer(html string, col int) [][]string {
	html = strings.Replace(html, "\n", "", -1)
	html = strings.Replace(html, "\r", "", -1)
	html = strings.Replace(html, "\t", "", -1)
	log.Println("html=========================================>", html)

	tdCompile := regexp.MustCompile(`<td>(.*?)</td>`)
	tds := tdCompile.FindAllString(html, -1)
	log.Println(tds)
	//根据tr标签个数确认有多少行数据
	trCompile := regexp.MustCompile(`<tr>(.*?)</tr>`)
	trs := trCompile.FindAllString(html, -1)
	row := len(trs) - 1
	if row == 0 {
		return nil
	}
	target := make([][]string, row)
	for i := range target {
		target[i] = make([]string, col)
	}
	//将数据写入数组中
	for i := 0; i < row; i++ {
		for j := 0; j < col; j++ {
			target[i][j] = tds[i*col+j][4 : len(tds[i*col+j])-5]
			log.Println(target[i][j])
		}
	}
	return target
}
