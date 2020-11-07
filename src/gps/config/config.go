package config

import (
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"goPanel/src/gps/common"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

var (
	Conf           *Config
	ConfigFilePath = common.GetCurrentDir() + "/src/gps/config/conf.yaml"
)

type Config struct {
	App      *AppConfig
	Database *Database
}

type AppConfig struct {
	Debug         bool   `yaml:"debug"`
	LogLevel      uint32 `yaml:"log_level"`
	LogOutputType uint32 `yaml:"log_output_type"`
	LogOutputFlag uint32 `yaml:"log_output_flag"`
	LogPath       string `yaml:"log_path"`
}

func init() {
	Conf = new(Config)

	loadYamlConfig()
	loadEnvConfig()

	new(Database).initialization()
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
		log.Fatal("Error loading .env file")
	}
}
