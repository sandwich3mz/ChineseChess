package config

import (
	"chesss/tools"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
)

var Conf = new(Config)

type Config struct {
	Database `mapstructure:"database" json:"database" yaml:"database"`
}

type Database struct {
	Driver       string `mapstructure:"driver" json:"driver" yaml:"driver"`
	Host         string `mapstructure:"host" json:"host" yaml:"host"`
	Port         int    `mapstrcture:"port" json:"port" yaml:"port"`
	Database     string `mapstructure:"database" json:"database" yaml:"database"`
	UserName     string `mapstructure:"username" json:"username" yaml:"username"`
	Password     string `mapstructure:"password" json:"password" yaml:"password"`
	Charset      string `mapstructure:"charset" json:"charset" yaml:"charset"`
	MaxIdleConns int    `mapstructure:"max_idle_conns" json:"max_idle_conns" yaml:"max_idle_conns"`
	MaxOpenConns int    `mapstructure:"max_open_conns" json:"max_open_conns" yaml:"max_open_conns"`
}

type App struct {
	Env     string `mapstructure:"env"`
	Port    string `mapstructure:"port"`
	AppName string `mapstructure:"app_name"`
	AppUrl  string `mapstructure:"app_url"`
}

func InitConfig() (err error) {

	// 默认配置文件路径
	configPath := tools.GetRootPath("/config.yaml")
	log.Printf("===> config path: %s", configPath)
	// 初始化配置文件
	viper.SetConfigFile(configPath)
	viper.WatchConfig()
	// 观察配置文件变动
	viper.OnConfigChange(func(in fsnotify.Event) {
		log.Printf("config file has changed")
		if err = viper.Unmarshal(&Conf); err != nil {
			log.Printf("failed at unmarshal config file after change, err: %v", err)
		}
	})
	// 将配置文件读入 viper
	if err = viper.ReadInConfig(); err != nil {
		log.Printf("failed at ReadInConfig, err: %v", err)
	}
	// 解析到变量中
	if err = viper.Unmarshal(&Conf); err != nil {
		log.Printf("failed at Unmarshal config file, err: %v", err)
	}
	// 从环境变量中覆盖配置
	//viper.AutomaticEnv()
	// 返回 nil 或错误
	return err
}
