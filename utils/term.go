package utils

import (
	"log"
	"strconv"
	"time"
)

//学期相关工具

// GetCurrentTerm 获取当前学期
func GetCurrentTerm() (term string) {
	date := time.Now()
	nowDate := date.Format("2006-01-02")
	year1date := date.Format("2006-01-02")[0:4]
	//明年
	year2date := date.AddDate(1, 0, 0).Format("2006-01-02")[0:4]
	//判断为上学期还是下学期，以8月份为分隔
	month, _ := strconv.Atoi(nowDate[5:7])
	log.Println(month)
	if month > 8 {
		//下半年为上学期
		term = year1date + "-" + year2date + "-" + "1"
	} else {
		//上学期
		term = year1date + "-" + year2date + "-" + "2"
	}
	return term
}
