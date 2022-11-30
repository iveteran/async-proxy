package receiver

import (
	"fmt"
	"os"

	"matrix.works/fmx-async-proxy/bootstrap"
	"matrix.works/fmx-async-proxy/conf"
)

func newApp(appName, appOwner, version, buildNo string) *bootstrap.FapBootstrapper {
	var appTokens = map[string]string{conf.Cfg.Server.AppId: conf.Cfg.Server.AppToken}
	app := bootstrap.New(
		appName,
		appOwner,
		version,
		buildNo,
		appTokens,
	)
	app.Bootstrap()
	return app
}

func Startup(appName, appOwner, version, buildNo string) {
	app := newApp(appName, appOwner, version, buildNo)

	SetupProxyHandlers()

	addr := fmt.Sprintf("%s:%d", conf.Cfg.Server.ListenAddress, conf.Cfg.Server.ListenPort)
	err := app.Serve(addr)
	if err != nil {
		conf.Logger.Printf("Start server failed: %s\n", err.Error())
		fmt.Printf("Start server failed: %s\n", err.Error())
		os.Exit(1)
	}
}
