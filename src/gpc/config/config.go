package config

import (
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"goPanel/src/common"
	"goPanel/src/constants"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

var (
	Conf              *Config
	GpcConfigFileName = common.GetCurrentDir() + constants.CONFIG_PATH + constants.GPC_CONFIG_FILENAME
	GpcPidFileName    = common.GetCurrentDir() + constants.GPC_PID_PATH + constants.PID_FILENAME
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
	if err := common.InitDir(
		common.GetCurrentDir()+constants.RUNTIME_PATH,
		common.GetCurrentDir()+constants.CONFIG_PATH,
		common.GetCurrentDir()+constants.GPC_PID_PATH,
	); err != nil {
		log.Panic(err)
	}

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
	if !common.DirOrFileByIsExists(GpcConfigFileName) {
		c, err := yaml.Marshal(map[string]interface{}{
			"app": DefaultAppConfig,
		})
		if err != nil {
			log.Panic(err)
			return
		}

		// 写配置文件
		err = ioutil.WriteFile(GpcConfigFileName, c, 0755)
		if err != nil {
			log.Panic("配置文件写入错误！", err)
		}
	}

	yamlFile, err := ioutil.ReadFile(GpcConfigFileName)
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
