package main

import (
	"os"
	"os/signal"
	"palo-alto/api"
	"palo-alto/config"
	dictionary "palo-alto/dictionary/file_dictionary"
	logger "palo-alto/logging"
	metric "palo-alto/metric/local_metric"
	"palo-alto/service"
	"syscall"
)

func main() {
	cfg, err := config.InitConf()

	if err != nil {
		logger.Log.Errorf("Cannot decode config: %s", err.Error())
		return
	}

	// Init metrics
	m := metric.New()
	go m.Start()
	defer m.Stop()

	// Init dictionary
	d := dictionary.New(cfg, m)
	go d.Start()
	defer d.Stop()
	err = d.Initialize()

	if err != nil {
		logger.Log.Errorf("Failed to init dictionary: %s", err.Error())
		return
	}

	// Init service
	srv := service.New(d)

	// api module
	a := api.New(cfg, m, srv)
	go a.Start()
	defer a.Stop()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
}
