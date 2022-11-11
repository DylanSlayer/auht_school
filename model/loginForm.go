package model

// LoginForm 用于登录的数据模型
type LoginForm struct {
	UserAccount string `json:"userAccount"`
	Encoded     string `json:"encoded"`
}
