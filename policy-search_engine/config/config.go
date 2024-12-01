package config

import (
	"fmt"
	"github.com/spf13/viper"
)

var V *viper.Viper

func Init() {
	V = viper.New()           //使用 viper.New() 创建了一个新的 viper 实例，并将其赋值给变量 V
	V.SetConfigName("config") //设置配置文件的名称为 "config"。这意味着 viper 会尝试加载名为 "config" 的配置文件（具体扩展名如 .json、.yaml、.toml 等，由配置文件的实际内容或默认设置决定）
	V.AddConfigPath("config") //添加一个配置文件的搜索路径 "config"。viper 会在这个路径下查找配置文件。如果 "config" 目录下有名为 "config" 的配置文件（如 "config/config.yaml"），那么 viper 就会加载它。
	err := V.ReadInConfig()   //调用 V.ReadInConfig() 方法来读取配置文件。如果配置文件存在且格式正确，那么配置信息就会被加载到 V 中。
	if err != nil {
		fmt.Printf("read config failed: %v\n", err)
	}
}
