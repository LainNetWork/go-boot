package application

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

// AppContext 应用全局公共配置
type AppContext struct {
	ConfigFilePrefix string //配置文件前缀，默认为application
	configs          []interface{}
	Config           struct {
		Boot struct {
			Active string `yaml:"active"` //当前激活环境，取application-{profile}.yaml中的配置
			Log    struct {
				Level      string `yaml:"level"`      //日志级别，默认info
				Path       string `yaml:"path"`       //日志位置，默认运行路径.
				SaveType   string `yaml:"saveType"`   //保存格式，json/text
				FileName   string `yaml:"fileName"`   //文件名
				MaxSize    int    `yaml:"maxSize"`    // 最大单文件大小（M）
				MaxBackups int    `yaml:"maxBackups"` // 最大备份数
				MaxAge     int    `yaml:"maxAge"`     // 单个文件记录最大时长（天）
				Compress   bool   `yaml:"compress"`   // 存档是否压缩
			}
		}
	}
}

func configDefault(v *viper.Viper) {
	v.SetDefault("boot.active", "dev")
	//日志
	v.SetDefault("boot.log.level", "info")
	v.SetDefault("boot.log.path", ".")
	v.SetDefault("boot.log.saveType", "text")
	v.SetDefault("boot.log.fileName", "boot")
	v.SetDefault("boot.log.maxSize", 10)
	v.SetDefault("boot.log.maxBackups", 3)
	v.SetDefault("boot.log.maxAge", 30)
	v.SetDefault("boot.log.compress", true)
}

var Context = &AppContext{
	ConfigFilePrefix: "application",
}

func (ctx *AppContext) RegisterConfig(config interface{}) {
	ctx.configs = append(ctx.configs, config)
}

func (ctx *AppContext) RegisterConfigs(config ...interface{}) {
	ctx.configs = append(ctx.configs, config...)
}

//Env 获取当前环境
func (ctx *AppContext) Env() string {
	return ctx.Config.Boot.Active
}

func (ctx *AppContext) Init() {
	//初始化环境、处理命令行传参
	ctx.prepare()
	//加载配置文件
	ctx.loadConfig()
	//根据配置，设置日志组件
	ctx.configLog()
}

func (ctx *AppContext) prepare() {
	//解析输入参数
	//反射处理注册的配置对象

	//args := os.Args
	//valueMap := make(map[string]string)
	//for _,arg := range args{
	//	if strings.HasPrefix(arg,"-") && strings.Contains(arg,"="){
	//		strings.SplitN()
	//	}
	//}
}

func (ctx *AppContext) configLog() {
	logConfig := ctx.Config.Boot.Log
	if logConfig.SaveType == "json" {
		log.SetFormatter(&log.JSONFormatter{})
	} else {
		log.SetFormatter(&log.TextFormatter{})
	}
	log.SetReportCaller(true)
	level, err := log.ParseLevel(logConfig.Level)
	if err == nil {
		log.SetLevel(level)
	} else {
		log.SetLevel(log.InfoLevel)
	}
	logFilePath := strings.TrimSuffix(strings.TrimSuffix(logConfig.Path, "/"), "\\") + "/" + logConfig.FileName + ".log"
	log.SetOutput(io.MultiWriter(&lumberjack.Logger{
		Filename:   logFilePath,
		MaxSize:    logConfig.MaxSize, // megabytes
		MaxBackups: logConfig.MaxBackups,
		MaxAge:     logConfig.MaxAge,   //days
		Compress:   logConfig.Compress, // disabled by default
	}, os.Stdout))
}

const defaultConfigKey = "default"

func (ctx *AppContext) loadConfig() {
	configMap := make(map[string]*viper.Viper)
	//解析配置文件。
	//判断有无指定配置文件，无配置文件则在缺省目录寻找默认配置文件
	var suffix = ".yaml"
	if files, err := ioutil.ReadDir("."); err == nil {
		for _, file := range files {
			configFileName := file.Name()
			if !file.IsDir() &&
				strings.HasPrefix(configFileName, ctx.ConfigFilePrefix) &&
				strings.HasSuffix(configFileName, suffix) {
				v := viper.New()
				v.SetConfigFile(configFileName)
				err := v.ReadInConfig()
				if err != nil {
					log.Fatal("读取配置文件失败！", err)
				}
				if configFileName == ctx.ConfigFilePrefix+suffix {
					//application.yaml 默认配置文件
					configMap[defaultConfigKey] = v
				} else {
					envTemp := strings.TrimSpace(strings.TrimPrefix(strings.TrimSuffix(configFileName, suffix), ctx.ConfigFilePrefix))
					if strings.HasPrefix(envTemp, "-") {
						env := strings.TrimPrefix(envTemp, "-")
						log.Infof("读取到环境 %s", env)
						configMap[env] = v
					}
				}
			}
		}
	}
	//先解析default
	vDefault := configMap[defaultConfigKey]
	if vDefault == nil {
		fmt.Print("未找到配置文件application.yaml,使用默认配置")
		vDefault = viper.New()
	}
	//设置配置缺省值
	configDefault(vDefault)
	err := vDefault.Unmarshal(&Context.Config)
	if err != nil {
		fmt.Print("解析配置文件出错！", err)
		os.Exit(1)
	}

	vEnv := configMap[Context.Config.Boot.Active]
	ctx.RegisterConfig(&Context.Config)
	for _, configStruct := range Context.configs {
		err := vDefault.Unmarshal(configStruct)
		if err != nil {
			fmt.Print("解析配置文件出错！", err)
			os.Exit(1)
		}
		if vEnv != nil {
			err := vEnv.Unmarshal(configStruct) //如启用了其他环境，则其他环境的配置覆盖默认配置
			if err != nil {
				fmt.Print("解析配置文件出错！", err)
				os.Exit(1)
			}
		}
	}
}

func (ctx *AppContext) WriteDefaultConfig() {
	v := viper.New()
	v.AddConfigPath(".")
	v.SetConfigName("application")
	v.SetConfigType("yaml")
	configDefault(v)
	err := v.SafeWriteConfig()
	if err != nil {
		log.Warning("錯誤 ", err)
	}
}
