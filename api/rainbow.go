package api

import (
	"auht_school/common"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

type Rainbow struct {
	Code   int            `json:"code"`
	Msg    string         `json:"msg"`
	Result RainbowContent `json:"result"`
}
type RainbowContent struct {
	Content string `json:"content"`
}

func ApiRainBow() (*Rainbow, error) {
	data := &Rainbow{
		Code:   201,
		Msg:    "请求失败",
		Result: RainbowContent{Content: "请求失败"},
	}
	resp, err := http.Get("https://apis.tianapi.com/caihongpi/index?key=" + common.Settings.Api.RainBow)
	if err != nil {
		return data, err
	}
	all, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return data, err
	}
	err = json.Unmarshal(all, data)
	if err != nil {
		return data, err
	}
	if data.Code == 150 {
		err = errors.New("api调用次数不足")
	}
	return data, err
}
