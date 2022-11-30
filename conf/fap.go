package conf

import (
	"log"

	"matrix.works/fmx-common/conf"
)

var Logger *log.Logger

var Cfg *MyConfig

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

type MyConfig struct {
	conf.FmxConfig
	Mq         MqInfo
	RouteTable map[string]string
	ReportApis map[string]StatusReportApi
}

func CreateGlobalConfig(filename string, logger *log.Logger) *MyConfig {
	Cfg = &MyConfig{
		FmxConfig: conf.FmxConfig{
			FilePath: filename,
		},
	}
	conf.Cfg = &Cfg.FmxConfig
	conf.LoadConfig(filename, Cfg, logger)
	conf.CheckRequiredOptions()
	return Cfg
}
