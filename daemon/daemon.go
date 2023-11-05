package daemon

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/jessevdk/go-flags"

	"matrix.works/async-proxy/bootstrap"
	"matrix.works/async-proxy/logger"
)

const (
	APP_OWNER   = "Matrixworks(ShenZhen) Information Technologies Co.,Ltd."
	APP_VERSION = "1.0.0"
)

type FnCleanup = func(*Daemon) error

type Daemon struct {
	AppName  string
	AppOwner string
	CmdOpts  *CommandOptions
	Logger   *log.Logger
	Ready    bool
}

func New(appName, appOwner, version, buildNo string, cleanup FnCleanup) *Daemon {
	d := &Daemon{AppName: appName, AppOwner: appOwner}
	d.CmdOpts = parseCommandLine(appName, appOwner, version, buildNo)
	d.Logger = logger.CreateFileLogger(appName, d.CmdOpts.LogFile)
	d.setupSignalHandler(cleanup)
	d.Ready = true
	return d
}

func (d *Daemon) Start() {

	if !d.Ready {
		panic("Not initialize the daemon")
	}

	select {}
}

func (d *Daemon) setupSignalHandler(cleanup FnCleanup) {
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT)
	go func() {
		for {
			<-s
			cleanup(d)
			println("Exit")
			os.Exit(0)
		}
	}()
}

type CommandOptions struct {
	bootstrap.CommandOptions
	LogFile string `short:"l" long:"logfile" description:"Log file name"`
}

func parseCommandLine(appName, appOwner, version, buildNo string) *CommandOptions {
	cmdOpts := new(CommandOptions)
	parser := flags.NewParser(cmdOpts, flags.Default)
	if _, err := parser.Parse(); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			fmt.Fprintf(os.Stderr, "[parseCommandLine]: %s\n", err.Error())
			os.Exit(1)
		}
	}

	if cmdOpts.VersionFlag {
		fmt.Printf("Fimatrix(%s) Version: %s, Build: %s, Copyright: %s\n",
			appName, version, buildNo, appOwner)
		os.Exit(0)
	}

	fmt.Printf("[parseCommandLine] Config file: %s\n", cmdOpts.ConfigPath)

	if _, err := os.Stat(cmdOpts.ConfigPath); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "[parseCommandLine]: %s\n", err.Error())
		os.Exit(1)
	}

	return cmdOpts
}
