package bootstrap

import (
	"log"
	"net/http"
)

type FapCommandOptions struct {
	*CommandOptions
	// add more options here
}

type FapBootstrapper struct {
	*Bootstrapper
	CmdOpts *FapCommandOptions
}

func New(
	appName, appOwner, appVersion, appBuildNo string,
	tokenTable map[string]string, cfgList ...Configurator,
) *FapBootstrapper {

	b := &FapBootstrapper{
		Bootstrapper: NewBootstrapper(
			appName,
			appOwner,
			appVersion,
			appBuildNo,
			tokenTable,
			cfgList...,
		),
		CmdOpts: &FapCommandOptions{
			&CommandOptions{},
		},
	}

	return b
}

func (this *FapBootstrapper) ParseCommandLine() {
	this.Bootstrapper.ParseCommandLine(this.CmdOpts.CommandOptions)
}

func (this *FapBootstrapper) Serve(addr string) error {
	server := &http.Server{
		Addr: addr,
	}

	log.Printf("Listening on %s...", addr)
	return server.ListenAndServe()
}
