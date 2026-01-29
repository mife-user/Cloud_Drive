package conf

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/viper"
)

// 获取全局配置
func GetConfig() *Config {
	return &globalConfig
}

// 加载配置文件
func LoadConfig() (*Config, error) {
	v := viper.New()
	//主要配置文件目录
	path := filepath.Join("config")
	v.AddConfigPath(path)
	//配置文件名称和类型
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	//环境变量配置
	v.AutomaticEnv()
	v.SetEnvPrefix("CLOUDPAN")
	v.AllowEmptyEnv(true)
	//读取主配置文件
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("加载主配置失败：%w", err)
	}
	//变更环境配置文件
	if env := v.GetString("env"); env != "" {
		v.SetConfigName("config." + env)
		if err := v.MergeInConfig(); err != nil {
			return nil, fmt.Errorf("加载 %s 配置失败：%w", env, err)
		}
	}
	//配置到结构体
	if err := v.Unmarshal(&globalConfig); err != nil {
		return nil, fmt.Errorf("解析配置失败：%w", err)
	}
	return &globalConfig, nil
}

// 检查配置
func StatusConfig() error {
	// MySQL配置检查
	if globalConfig.Mysql.Dsn == "" {
		return fmt.Errorf("mysql连接未配置")
	}
	// Redis配置检查
	if globalConfig.Redis.Port == "" {
		return fmt.Errorf("redis端口未配置")
	}
	return nil
}
