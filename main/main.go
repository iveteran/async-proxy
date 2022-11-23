package main

import (
	"matrix.works/fmx-async-proxy/conf"
	"matrix.works/fmx-async-proxy/receiver"
	"matrix.works/fmx-async-proxy/sender"
	"matrix.works/fmx-common/daemon"
)

const (
	APP_NAME  = "fap"
	APP_OWNER = "Matrixworks(ShenZhen) Information Technologies Co.,Ltd."
)

var (
	Version string = "Unknown"
	BuildNo string = "Unknown"
)

func cleanup(d *daemon.FmxDaemon) error {
	sender.Cleanup()
	return nil
}

func main() {
	daemon := daemon.NewFmxDaemon(APP_NAME, APP_OWNER, Version, BuildNo, cleanup)
	conf.CreateGlobalConfig(daemon.CmdOpts.ConfigPath, daemon.Logger)
	conf.Logger = daemon.Logger

	go receiver.Startup(APP_NAME, APP_OWNER, Version, BuildNo)
	//go sender.Startup(APP_NAME, APP_OWNER, Version, BuildNo)
	sender.Startup(APP_NAME, APP_OWNER, Version, BuildNo)

	daemon.Start()
}
