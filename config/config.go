package config

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// AppConfig 项目配置
type AppConfig struct {
	Version string `mapstructure:"version"` // 版本号
}

// Server 服务器
type Server struct {
	Ip      string `mapstructure:"ip"`      // 运行IP
	Port    uint   `mapstructure:"port"`    // 运行端口
	Timeout uint   `mapstructure:"timeout"` // 超时时间
}

// Img 图片
type Img struct {
	Dir string `mapstructure:"dir"` // 存放目录
}

// Config 总配置
type Config struct {
	App  *AppConfig `mapstructure:"app"`
	Http *Server    `mapstructure:"http"`
	Img  *Img       `mapstructure:"Img"`
}

// splitConfigFile 通过文件路径获取目录、文件名、扩展名
func splitConfigFile(configFile string) (dir string, fileName string, extName string, err error) {
	if len(configFile) == 0 {
		err = errors.New(configFile + " is empty")
		return
	}
	configFiles := strings.Split(configFile, "/")
	lens := len(configFiles) - 1
	if lens == 0 {
		dir = "."
	} else {
		dir = strings.Join(configFiles[:lens], "/")
	}
	files := strings.Split(configFiles[lens], ".")
	if len(files) <= 1 {
		err = errors.New(configFile + " file name is empty")
		return
	}
	fileName = files[0]
	extName = files[1]
	return
}

// InitViper 初始化配置文件
// configFile 配置文件
// isEnv      是否获取环境变量环境
// loadData   装载的数据结构指针类型
func InitViper(configFile string, isEnv bool, loadData interface{}) error {
	dir, fileName, extName, err := splitConfigFile(configFile)
	if err != nil {
		return err
	}

	v := viper.New()
	// 设置配置文件的名字
	v.SetConfigName(fileName)
	// 添加配置文件所在的路径
	v.AddConfigPath(dir)
	// 设置配置文件类型
	v.SetConfigType(extName)
	if err = v.ReadInConfig(); err != nil {
		return err
	}
	if err = v.Unmarshal(loadData); err != nil {
		return err
	}
	return nil
}

var conf *Config

// Init 初始化配置
func Init(configFile string) {
	conf = &Config{
		App:  &AppConfig{},
		Http: &Server{},
	}
	if err := InitViper(configFile, true, conf); err != nil {
		panic(err)
	}
}

// Get Config
func Get() *Config {
	return conf
}

func HttpAddr() string {
	if conf.Http.Port == 0 {
		conf.Http.Port = 8080
	}
	return fmt.Sprintf("%s:%d", conf.Http.Ip, conf.Http.Port)
}

func ImgDir() string {
	if conf.Img.Dir == "" {
		conf.Img.Dir = "img/"
	}
	return conf.Img.Dir
}
