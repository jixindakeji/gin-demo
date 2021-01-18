package setting

import (
	"flag"
	"fmt"
	"github.com/spf13/viper"
	"os"
	"time"
)

type App struct {
	Salt         string `mapstructure:"salt" json:"-" yaml:"salt"`
	PageSize     int    `mapstructure:"page-size" json:"PageSize" yaml:"page-size"`
	JwtTokenName string `mapstructure:"jwt-token-name" json:"JwtTokenName" yaml:"jwt-token-name"`
	JwtSecret    string `mapstructure:"jwt-secret" json:"JwtSecret" yaml:"jwt-secret"`
	AdminSecret  string `mapstructure:"admin-secret" json:"-" yaml:"jwt-secret"`
	LogLevel     string `mapstructure:"log-level" json:"LogLevel" yaml:"log-level"`
	LogType      string `mapstructure:"log-type" json:"LogType" yaml:"log-type"`
	LogPath      string `mapstructure:"log-path" json:"LogPath" yaml:"log-path"`
	LogName      string `mapstructure:"log-name" json:"LogName" yaml:"log-name"`
}
type Server struct {
	Port         int           `mapstructure:"port" json:"port" yaml:"port"`
	ReadTimeout  time.Duration `mapstructure:"read-timeout" json:"ReadTimeout" yaml:"read-timeout"`
	WriteTimeout time.Duration `mapstructure:"write-timeout" json:"WriteTimeout" yaml:"write-timeout"`
}

type DataBase struct {
	Type        string `mapstructure:"type" json:"type" yaml:"type"`
	Host        string `mapstructure:"host" json:"host" yaml:"host"`
	User        string `mapstructure:"user" json:"user" yaml:"user"`
	Password    string `mapstructure:"password" json:"password" yaml:"password"`
	Name        string `mapstructure:"name" json:"name" yaml:"name"`
	TablePrefix string `mapstructure:"table-prefix" json:"TablePrefix" yaml:"table-prefix"`
}

type Config struct {
	RunMode  string `mapstructure:"run-mode" json:"RunMode" yaml:"run-mode"`
	APP      *App
	Server   *Server
	DataBase *DataBase
}

const (
	ConfigEnv  = "APP_CONFIG"
	ConfigFile = "conf/app.yaml"
)

var GConfig Config

func init() {
	var config string

	flag.StringVar(&config, "c", "", "choose config file.")
	flag.Parse()
	if config == "" { // 优先级: 命令行 > 环境变量 > 默认值
		if configEnv := os.Getenv(ConfigEnv); configEnv == "" {
			config = ConfigFile
			fmt.Printf("使用config的默认值,config的路径为%v\n", ConfigFile)
		} else {
			config = configEnv
			fmt.Printf("使用APP_CONFIG环境变量,config的路径为%v\n", config)
		}
	} else {
		fmt.Printf("使用命令行的-c参数传递的值,config的路径为%v\n", config)
	}

	v := viper.New()
	v.SetConfigFile(config)
	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("解析配置文件失败: %s \n", err))
	}
	/*
		v.WatchConfig()
		v.OnConfigChange(func(e fsnotify.Event) {
			fmt.Println("配置文件改变:", e.Name)
			if err := v.Unmarshal(&GConfig); err != nil {
				fmt.Println(err)
			}
		})
	*/
	if err := v.Unmarshal(&GConfig); err != nil {
		fmt.Println(err)
	}
}
