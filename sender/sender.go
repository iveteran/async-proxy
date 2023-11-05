package sender

import (
	"fmt"
	"os"

	"matrix.works/async-proxy/conf"
	"matrix.works/async-proxy/logger"
	"matrix.works/async-proxy/mq"
)

var (
	cc *mq.MqConsumerCreator
)

func Cleanup() error {
	cc.Cleanup()
	return nil
}

func Startup(appName, appOwner, version, buildNo string) {
	cc = &mq.MqConsumerCreator{}
	err := cc.Init(conf.Cfg.Mq.ConnectionName)
	if err != nil {
		logger.Logger.Printf("Start consumer failed: %s\n", err.Error())
		fmt.Printf("Start consumer failed: %s\n", err.Error())
		os.Exit(1)
	}
}
