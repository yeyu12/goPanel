package config

import (
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"goPanel/src/common"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

var (
	Conf              *Config
	GpcConfigFilePath = common.GetCurrentDir() + "/config/gpc.yaml"
)

type Config struct {
	App *AppConfig
	Ssh *SshConfig
}

type AppConfig struct {
	Debug                bool   `yaml:"debug"`
	LogLevel             uint32 `yaml:"log_level"`
	LogOutputType        uint32 `yaml:"log_output_type"`
	LogOutputFlag        uint32 `yaml:"log_output_flag"`
	LogPath              string `yaml:"log_path"`
	ServerHost           string `yaml:"server_host"`
	ServerPort           string `yaml:"server_port"`
	LocalName            string `yaml:"local_name"`
	ControlHeartbeatTime int64  `yaml:"control_heartbeat_time"`
	ControlReconnTcpTime int64  `yaml:"control_reconn_tcp_time"`
	UidPath              string `yaml:"uid_path"`
	UserDir              string
	Uid                  string
}

var DefaultAppConfig = map[string]interface{}{
	"debug":                   false,
	"log_level":               2,
	"log_output_type":         0,
	"log_output_flag":         1,
	"log_path":                "./runtime/log/",
	"server_host":             "192.168.28.124",
	"server_port":             10010,
	"local_name":              "主机",
	"control_heartbeat_time":  60,
	"control_reconn_tcp_time": 5,
	"uid_path":                "./runtime/uid/",
}

func init() {
	Conf = new(Config)

	loadYamlConfig()
	loadEnvConfig()

	userDir, err := common.UserDir()
	if err != nil {
		log.Panic("获取用户目录失败！", err)
	}
	Conf.App.UserDir = userDir

	new(SshConfig).initialization()
}

func loadYamlConfig() {
	if !common.DirOrFileByIsExists(GpcConfigFilePath) {
		c, err := yaml.Marshal(map[string]interface{}{
			"app": DefaultAppConfig,
		})
		if err != nil {
			log.Panic(err)
			return
		}

		// 写配置文件
		err = ioutil.WriteFile(GpcConfigFilePath, c, 0755)
		if err != nil {
			log.Panic("配置文件写入错误！", err)
		}
	}

	yamlFile, err := ioutil.ReadFile(GpcConfigFilePath)
	if err != nil {
		log.Panic("yamlFile.Get err #", err)
	}

	if err = yaml.Unmarshal(yamlFile, Conf); err != nil {
		log.Panic("Unmarshal: %v", err)
	}
}

func loadEnvConfig() {
	_ = godotenv.Load()
}
