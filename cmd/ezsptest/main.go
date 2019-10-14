package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/conthing/ezsp/ash"
	"github.com/conthing/ezsp/ezsp"
	"github.com/conthing/ezsp/zgb"
	"github.com/conthing/utils/common"
)

type stConfig struct {
	Serial stSerialConfig `json:"serial"`
}

type stSerialConfig struct {
	Name   string
	Baud   uint
	RtsCts bool
}

var cfg = stConfig{}

func boot(_ interface{}) (needRetry bool, err error) {
	var cfgfile string
	var action string

	//解析命令行参数 -c <cfgfile> 默认configuration.toml
	flag.StringVar(&cfgfile, "config", "configuration.toml", "Specify a config file other than default.")
	flag.StringVar(&cfgfile, "c", "configuration.toml", "Specify a config file other than default.")
	flag.StringVar(&action, "act", "", "Specify a init action form/remove for form or remove network")
	flag.Parse()

	common.InitLogger(&common.LoggerConfig{Level: "DEBUG", SkipCaller: true})

	err = common.LoadConfig(cfgfile, &cfg)
	if err != nil {
		common.Log.Errorf("failed to load config %v", err)
		return false, err
	}

	common.Log.Infof("load cfg success %+v", cfg)

	err = ash.AshSerialOpen(cfg.Serial.Name, cfg.Serial.Baud, cfg.Serial.RtsCts)
	if err != nil {
		common.Log.Errorf("failed to open serial %v", cfg.Serial.Name)
		return false, err
	}

	// Time it took to start service
	common.Log.Infof("Open Serial success port=%s baud=%d", cfg.Serial.Name, cfg.Serial.Baud)

	return false, nil
}

func main() {
	if common.Bootstrap(boot, nil, 1000, 500) != nil {
		return
	}

	// 1.SIGINT 2.httpserver 3.serial recv
	errs := make(chan error, 3)

	listenForInterrupt(errs)
	ash.AshStartTransceiver(ezsp.AshRecvImp, errs)
	//startHTTPServer(errs, cfg.HTTP.Port)
	startTickRunning(errs)

	// recv error channel
	c := <-errs
	common.Log.Errorf("terminating: %v", c)
	os.Exit(0)
}

func listenForInterrupt(errChan chan error) {
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errChan <- fmt.Errorf("%s", <-c)
	}()
}

func startTickRunning(errChan chan error) {
	go zgb.TickRunning(errChan)
}
