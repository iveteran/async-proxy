package bootstrap

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jessevdk/go-flags"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"

	mylog "matrix.works/async-proxy/logger"
)

var (
	SysTimeLocation  *time.Location
	SysTimeForm      string
	SysTimeFormShort string
)

type Configurator func(bootstrapper *Bootstrapper)

type CommandOptions struct {
	ConfigPath  string `short:"c" long:"config" description:"Configure file path"`
	VersionFlag bool   `short:"v" long:"version" description:"Show version"`
}

type Bootstrapper struct {
	*iris.Application
	AppName      string
	AppOwner     string
	AppVersion   string
	AppBuildNo   string
	AppSpawnDate time.Time
	Logger       *log.Logger
	AppTokens    map[string]string
}

func NewBootstrapper(appName, appOwner, appVersion, appBuildNo string,
	tokenTable map[string]string, cfgList ...Configurator,
) *Bootstrapper {

	filePath := fmt.Sprintf("./%s.log", appName)
	logger := mylog.CreateFileLogger(appName, filePath) // TODO: use iris logger
	b := &Bootstrapper{
		Application:  iris.New(),
		AppName:      appName,
		AppOwner:     appOwner,
		AppVersion:   appVersion,
		AppBuildNo:   appBuildNo,
		AppSpawnDate: time.Now(),
		Logger:       logger,
		AppTokens:    tokenTable,
	}

	for _, cfg := range cfgList {
		cfg(b)
	}

	return b
}

func (this *Bootstrapper) Bootstrap() *Bootstrapper {
	this.SetupErrorHandler()

	this.SetupCron()

	this.Use(recover.New())
	this.Use(logger.New())

	this.SetupSignalHandler()

	return this
}

func (this *Bootstrapper) Listen(addr string, cfgList ...iris.Configurator) error {
	err := this.Run(iris.Addr(addr), cfgList...)

	if err != nil {
		log.Fatal("Bootstrap.Listen error ", err)
	}
	return err
}

func (this *Bootstrapper) SetupErrorHandler() {
	this.OnAnyErrorCode(func(ctx iris.Context) {
		err := iris.Map{
			"app":     this.AppName,
			"status":  ctx.GetStatusCode(),
			"message": ctx.Values().GetString("message"),
		}

		if jsonOutput := ctx.URLParamExists("json"); jsonOutput {
			ctx.JSON(err)
			return
		}

		ctx.JSON(err)
	})
}

func (this *Bootstrapper) Configure(cfgList ...Configurator) {
	for _, cfg := range cfgList {
		cfg(this)
	}
}

func (this *Bootstrapper) ParseCommandLine(cmdOpts *CommandOptions) {
	parser := flags.NewParser(cmdOpts, flags.Default)
	if _, err := parser.Parse(); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			log.Printf("Parse command line: %s\n", err.Error())
			os.Exit(1)
		}
	}

	if cmdOpts.VersionFlag {
		fmt.Printf("Fimatrix(%s) Version: %s, BuildNo: %s, Copyright: %s\n",
			this.AppName, this.AppVersion, this.AppBuildNo, this.AppOwner)
		os.Exit(0)
	}

	log.Printf("Config file: %s\n", cmdOpts.ConfigPath)

	if _, err := os.Stat(cmdOpts.ConfigPath); os.IsNotExist(err) {
		log.Fatal(err)
	}
}

func (this *Bootstrapper) SetupSignalHandler() {
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGUSR1)
	go func() {
		for {
			<-s
			//conf.Cfg.ReloadConfig(this.Logger) TODO: uncomment this line
			log.Println("Reloaded config")
		}
	}()
}

func (this *Bootstrapper) SetupCron() {
	// TODO
}

func (this *Bootstrapper) GetAppToken(appId string) string {
	return this.AppTokens[appId]
}
