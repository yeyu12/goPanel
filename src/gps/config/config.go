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
	GpsConfigFilePath = common.GetCurrentDir() + "/config/gps.yaml"
)

type Config struct {
	App      *AppConfig
	Database *Database
}

type AppConfig struct {
	Debug          bool   `yaml:"debug"`
	LogLevel       uint32 `yaml:"log_level"`
	LogOutputType  uint32 `yaml:"log_output_type"`
	LogOutputFlag  uint32 `yaml:"log_output_flag"`
	LogPath        string `yaml:"log_path"`
	HttpPort       int    `yaml:"http_port"`
	ControlPort    int    `yaml:"control_port"`
	RelayStartPort int    `yaml:"relay_start_port"`
}

// 默认配置
var DefaultConfigApp = map[string]interface{}{
	"debug":            false,
	"log_level":        2,
	"log_output_type":  0,
	"log_output_flag":  0,
	"log_path":         "./runtime/log/",
	"http_port":        10000,
	"control_port":     10010,
	"relay_start_port": 10086,
}

func init() {
	Conf = new(Config)

	loadYamlConfig()
	loadEnvConfig()

	new(Database).initialization()
}

func loadYamlConfig() {
	if !common.DirOrFileByIsExists(GpsConfigFilePath) {
		c, err := yaml.Marshal(map[string]interface{}{
			"app": DefaultConfigApp,
		})
		if err != nil {
			log.Panic(err)
			return
		}

		// 写配置文件
		err = ioutil.WriteFile(GpsConfigFilePath, c, 0755)
		if err != nil {
			log.Panic("配置文件写入错误！", err)
		}
	}

	yamlFile, err := ioutil.ReadFile(GpsConfigFilePath)
	if err != nil {
		log.Panic("yamlFile.Get err #%v ", err)
	}

	if err = yaml.Unmarshal(yamlFile, Conf); err != nil {
		log.Panic("Unmarshal: %v", err)
	}
}

func loadEnvConfig() {
	_ = godotenv.Load()
}
