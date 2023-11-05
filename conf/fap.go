package conf

import (
	"log"
)

var Cfg *FapConfig

type StatusReportApi struct {
	Url    string
	Method string
	Args   string
}

type QueueInfo struct {
	Name         string
	UnackedLimit int64
	NumConsumers int64
	MessageTTL   int64
	Backend      string
}

type MqInfo struct {
	ConnectionName string
	DefaultQueue   QueueInfo
	TopicQueues    map[string]QueueInfo
}

type BackendInfo struct {
	Host     string
	AppId    string
	AppToken string
}

type Backends struct {
	AppId    string
	AppToken string
	Backends map[string]BackendInfo
}

type FapConfig struct {
	Config
	Mq         MqInfo
	RouteTable map[string]string
	ReportApis map[string]StatusReportApi
}

func CreateGlobalConfig(filename string, logger *log.Logger) *FapConfig {
	Cfg = &FapConfig{
		Config: Config{
			FilePath: filename,
		},
	}
	BaseCfg = &Cfg.Config
	LoadConfig(filename, Cfg, logger)
	CheckRequiredOptions()
	return Cfg
}
