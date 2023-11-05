package main

import (
	"matrix.works/async-proxy/conf"
	"matrix.works/async-proxy/daemon"
	"matrix.works/async-proxy/logger"
	"matrix.works/async-proxy/receiver"
	"matrix.works/async-proxy/sender"
)

const (
	APP_NAME  = "fap"
	APP_OWNER = "Matrixworks(ShenZhen) Information Technologies Co.,Ltd."
)

var (
	Version string = "Unknown"
	BuildNo string = "Unknown"
)

func cleanup(d *daemon.Daemon) error {
	sender.Cleanup()
	return nil
}

func main() {
	daemon := daemon.New(APP_NAME, APP_OWNER, Version, BuildNo, cleanup)
	conf.CreateGlobalConfig(daemon.CmdOpts.ConfigPath, daemon.Logger)
	logger.Logger = daemon.Logger

	go receiver.Startup(APP_NAME, APP_OWNER, Version, BuildNo)
	//go sender.Startup(APP_NAME, APP_OWNER, Version, BuildNo)
	sender.Startup(APP_NAME, APP_OWNER, Version, BuildNo)

	daemon.Start()
}
