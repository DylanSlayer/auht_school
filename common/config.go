package common

import "github.com/spf13/viper"

var (
	Settings GlobalConfig
)

type GlobalConfig struct {
	ServInfo ServerConfig `mapstructure:"server"`
	//FakeDataInfo FakeDataConfig `mapstructure:"fake_data"`
	Api ApiConfig `mapstructure:"api"`
}
type ApiConfig struct {
	RainBow string `mapstructure:"rainbow-key"`
}

type ServerConfig struct {
	Port    int    `mapstructure:"port"`
	BaseUrl string `mapstructure:"base-url"`
}

////数据集
//type FakeDataConfig struct {
//	StudentNames []string `mapstructure:"student_names"`
//	Academies    []string `mapstructure:"academies"`
//	TeacherNames []string `mapstructure:"teacher_names"`
//	Addresses    []string `mapstructure:"addresses"`
//	Reason       []string `mapstructure:"reason"`
//	Approver     []string `mapstructure:"teacher_names"`
//	IdBegin      int      `mapstructure:"idBegin"`
//	IdEnd        int      `mapstructure:"idEnd"`
//}

// InitConfig 初始化配置
func InitConfig() {
	v := viper.New()
	v.SetConfigFile("./config.yml")
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}
	servconfig := GlobalConfig{}
	if err := v.Unmarshal(&servconfig); err != nil {
		panic(err)
	}
	Settings = servconfig
}
