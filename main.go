package main

import (
	"os"
	"os/signal"
	"palo-alto/api"
	"palo-alto/conf"
	"palo-alto/dictionary/fileDictionary"
	logger "palo-alto/logging"
	"palo-alto/metric"
	"palo-alto/service"
	"syscall"
)

// Module base module interface
type Module interface {
	Start()
	Stop()
	Title() string
}

func main() {
	cfg, err := conf.InitConf()

	if err != nil {
		logger.Log.Errorf("Cannot decode config: %s", err.Error())
		return
	}

	// Init metrics
	m := metric.New()

	// Init dictionary
	d := fileDictionary.New(cfg, m)
	err = d.Initialize()

	if err != nil {
		logger.Log.Errorf("Failed to init dictionary: %s", err.Error())
		return
	}

	// Init service
	srv := service.New(d)

	// api module
	apiModule := api.New(cfg, m, srv)

	RunModules(apiModule)
}

// RunModules runs each of the modules in a separate goroutine.
func RunModules(modules ...Module) {
	defer func() {
		for _, m := range modules {
			logger.Log.Infof("Stopping module %s", m.Title())
			m.Stop()
		}
		logger.Log.Infof("Stopped all modules")
	}()

	for _, m := range modules {
		logger.Log.Infof("Starting module %s", m.Title())
		go m.Start()
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
}
