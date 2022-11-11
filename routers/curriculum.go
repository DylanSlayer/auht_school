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
	"time"
)

//与课程表相关的router

func SetCurriculum(r *gin.Engine) {
	r.POST("/curriculum", getCurrirulum)
}
func getCurrirulum(c *gin.Context) {
	targetUrl := "http://jwxt.ahut.edu.cn/jsxsd/framework/main_index_loadkb.jsp"
	method := "POST"
	contentType := "application/x-www-form-urlencoded; charset=UTF-8"
	sjms := "9697EC28830EB8C8E0531583640A461C"
	type dateJson struct {
		Date string `json:"date"`
	}
	date := dateJson{}
	err := c.ShouldBind(&date)
	if err != nil {
		c.String(http.StatusBadRequest, "时间数据绑定失败")
		return
	}
	if date.Date == "all" {
		targetUrl = "http://jwxt.ahut.edu.cn/jsxsd/xskb/xskb_list.do"
		method = "GET"
	}
	if date.Date == "" {
		date.Date = time.Now().Format("2006-01-02")
	}

	log.Println("日期：", date.Date)
	formValues := url.Values{}
	formValues.Set("rq", date.Date)
	formValues.Set("sjmsValue", sjms)
	formDataStr := formValues.Encode()
	formDataBytes := []byte(formDataStr)
	formBytesReader := bytes.NewReader(formDataBytes)
	request, err := utils.GetRequest(c, method, targetUrl, contentType, formBytesReader)
	if err != nil {
		c.String(200, "请求失败")
		return
	}
	resp, _ := utils.Client.Do(request)
	defer resp.Body.Close()
	html_bytes, err := ioutil.ReadAll(resp.Body)
	html := string(html_bytes)
	html = strings.Replace(html, "\n", "", -1)
	html = strings.Replace(html, "\r", "", -1)
	html = strings.Replace(html, "\t", "", -1)
	//log.Println("html===>", html)
	var data []model.Curriculum

	if date.Date == "all" {

		//学期理论课表的处理
		//第几节课，最后一个是备注、注
		th := regexp.MustCompile(`<th width="70" height="28" align="center">(.*?)</th>`)
		//具体课程，最后一个是备注描述
		td := regexp.MustCompile(`<td(.*?)</td>`)
		match_th := th.FindAllString(html, -1)
		match_td := td.FindAllString(html, -1)
		if len(match_th) == 0 || len(match_td) == 0 {
			//无法查询，可能是登录失败或者登录过时
			c.JSON(200, gin.H{
				"code": 433,
				"msg":  "无法请求到数据，请稍后重试或联系管理员处理",
			})
		}
		log.Println("len td:", len(match_td))
		log.Println("len th:", len(match_th))
		for i, s := range match_th {
			//处理节次数据，7个
			if i != 0 {
				data = append(data, model.Curriculum{Class: s})
			}
		}
		//处理其余数据
		for i := 0; i < 7; i++ {
			if i < 6 {
				data[i].One = match_td[i*7+0]
			} else {
				data[i].One = match_td[len(match_td)-1]
			}
			if i < 6 {
				data[i].Two = match_td[i*7+1]
				data[i].Three = match_td[i*7+2]
				data[i].Four = match_td[i*7+3]
				data[i].Five = match_td[i*7+4]
				data[i].Six = match_td[i*7+5]
				data[i].Seven = match_td[i*7+6]
			}

		}
	} else {
		var rowData [48]string

		td := regexp.MustCompile(`<td.*?</td>`)
		match := td.FindAllString(html, -1)
		//log.Println(len(match))
		for i, s := range match {
			//log.Println(s)
			rowData[i] = s
			rowData[i] = strings.Replace(rowData[i], `<br/>`, "\n", -1)
			log.Println("rowData:====>", rowData[i])
		}
		//截取每一行的数据
		for i := 0; i < 6; i++ {
			data = append(data, model.Curriculum{
				Class: rowData[i*8+0],
				One:   rowData[i*8+1],
				Two:   rowData[i*8+2],
				Three: rowData[i*8+3],
				Four:  rowData[i*8+4],
				Five:  rowData[i*8+5],
				Six:   rowData[i*8+6],
				Seven: rowData[i*8+7],
			})
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": data,
		"date": date,
	})
}
