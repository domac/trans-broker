package main

import (
	"flag"
	"github.com/domac/rp"
	"github.com/domac/trans-broker/app"
	"github.com/domac/trans-broker/event"
	"github.com/domac/trans-broker/logger"
	"os"
)

const (
	ver = "v0.1"
)

var (
	flagSet     = flag.NewFlagSet("trans-brolers", flag.ExitOnError)
	configFile  = flagSet.String("config", "config.yml", "config file path")
	showVersion = flagSet.Bool("v", false, "show version")
	debug       = flagSet.Bool("debug", false, "open debug mode")
)

func init() {
	logger.InitLogger("/tmp/tb.log")
}

func main() {

	flagSet.Parse(os.Args[1:])

	app.SetVersion(ver)

	if *showVersion {
		app.ShowVersion()
		os.Exit(0)
	}

	appConfig, err := app.ParseConfigFile(*configFile)
	if err != nil {
		logger.Log().Error(err)
		panic(err)
	}

	//启动API服务
	go app.RunHTTPServer(appConfig)

	if *debug {
		rp.DEBUG_PROFILE()
	}

	//等待事件
	event.OnEvent(event.Event_EXIT, func() {
		app.ShutdownHTTPServer()
	})
	event.WaitEvent()
	event.EmitEvent(event.Event_EXIT)
}
