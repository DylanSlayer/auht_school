package utils

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"time"
)

type Cookie struct {
	Name       string        `json:"Name"`
	Value      string        `json:"Value"`
	Path       string        `json:"Path"`
	Domain     string        `json:"Domain"`
	Expires    time.Time     `json:"Expires"`
	RawExpires string        `json:"RawExpires"`
	MaxAge     int           `json:"MaxAge"`
	Secure     bool          `json:"Secure"`
	HttpOnly   bool          `json:"HttpOnly"`
	SameSite   http.SameSite `json:"SameSite"`
	Raw        string        `json:"Raw"`
}

var (
	Client *http.Client
)

func InitClient() {
	Client = &http.Client{
		Transport:     nil,
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       99999999999992,
	}
}

// GetRequest  获取一个请求
func GetRequest(c *gin.Context, method, url, contentType string, reader io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, reader)
	if err != nil {
		return &http.Request{}, err
	}
	//设置请求头
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	//req.Header.Set("Host", "jwxt.ahut.edu.cn")
	//req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Encoding", "")
	//req.Header.Set("Connection", "keep-alive")
	//req.Header.Set("Cache-Control", "no-cache")
	//req.Header.Set("Pragma", "no-cache")
	//req.Header.Set("Upgrade-Insecure-Requests", "1")
	//req.Header.Set("Origin", "http://jwxt.ahut.edu.cn")
	//获取暂存在客户端的cookie
	auth := c.Request.Header.Get("Authorization")
	if auth != "undefined" && auth != "" {
		req.Header.Set("Referer", "http://jwxt.ahut.edu.cn/jsxsd/framework/xsMain.jsp")
		var cookies []Cookie
		err = json.Unmarshal([]byte(auth), &cookies)
		for _, s := range cookies {
			cookie := &http.Cookie{
				Domain:   "jwxt.ahut.edu.cn",
				Name:     s.Name,
				Value:    s.Value,
				Path:     s.Path,
				HttpOnly: s.HttpOnly,
			}
			log.Println("要加入的cookie:", cookie)
			req.AddCookie(cookie)
		}
	}
	log.Println("req cookie:", req.Header.Get("Cookie"))
	return req, err
}
