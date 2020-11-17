package config

import (
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"goPanel/src/common"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

var (
	Conf           *Config
	ConfigFilePath = common.GetCurrentDir() + "/src/gpc/config/conf.yaml"
)

type Config struct {
	App *AppConfig
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
}

func init() {
	Conf = new(Config)

	loadYamlConfig()
	loadEnvConfig()

}

func loadYamlConfig() {
	yamlFile, err := ioutil.ReadFile(ConfigFilePath)
	if err != nil {
		log.Panic("yamlFile.Get err #%v ", err)
	}

	if err = yaml.Unmarshal(yamlFile, Conf); err != nil {
		log.Panic("Unmarshal: %v", err)
	}
}

func loadEnvConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Error("Error loading .env file")
	}
}
