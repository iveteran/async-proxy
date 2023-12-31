package receiver

import (
	"fmt"
	"os"

	"matrix.works/async-proxy/bootstrap"
	"matrix.works/async-proxy/conf"
	"matrix.works/async-proxy/logger"
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

	routeMap := conf.Cfg.RouteTable
	fmt.Printf("route table: %+v\n", routeMap)
	SetupProxyHandlers(routeMap)

	addr := fmt.Sprintf("%s:%d", conf.Cfg.Server.ListenAddress, conf.Cfg.Server.ListenPort)
	err := app.Serve(addr)
	if err != nil {
		logger.Logger.Printf("Start receiver failed: %s\n", err.Error())
		fmt.Printf("Start receiver failed: %s\n", err.Error())
		os.Exit(1)
	}
}
