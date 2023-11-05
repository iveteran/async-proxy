package conf

import (
	"log"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

var BaseCfg *Config

type ServerConfig struct {
	Host          string
	ListenAddress string
	ListenPort    int
	LogFile       string
	SysTimeForm   string
	SysDateForm   string
	SysTimeZone   string
	CookieSecret  string
	SignSecret    string
	AppId         string
	AppToken      string
}

type RedisConfig struct {
	Host     string
	Port     int
	Database int
}

type Config struct {
	Debug    bool
	FilePath string
	Server   ServerConfig
	Redises  map[string]RedisConfig
}

func LoadConfig(filePath string, cfg interface{}, logger *log.Logger) {
	filePath, err := filepath.Abs(filePath)
	if err != nil {
		panic(err)
	}
	if _, err := toml.DecodeFile(filePath, cfg); err != nil {
		panic(err)
	}
	if logger != nil {
		logger.Printf("Parse toml file once. filePath: %s\n", filePath)
		logger.Printf("Config: %+v\n", cfg)
	} else {
		log.Printf("Parse toml file once. filePath: %s\n", filePath)
		log.Printf("Config: %+v\n", cfg)
	}
}

func CheckRequiredOptions() {
	if Cfg == nil {
		log.Fatal("Not load configure")
	}
	if Cfg.Server.ListenAddress == "" || Cfg.Server.ListenPort == 0 {
		log.Fatalf("Invalid listen address or port: %s:%d", Cfg.Server.ListenAddress, Cfg.Server.ListenPort)
	}
}
