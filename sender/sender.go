package sender

import (
	"os"

	"matrix.works/fmx-async-proxy/conf"
	"matrix.works/fmx-async-proxy/mq"
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
		println(err.Error())
		os.Exit(1)
	}
}
