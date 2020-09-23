package main

import (
	"fmt"
	"github.com/shaolinjehzu/testGo/config"
	"github.com/shaolinjehzu/testGo/klines"
	"github.com/shaolinjehzu/testGo/monitoring"
	"github.com/shaolinjehzu/testGo/telegram"
	"github.com/shaolinjehzu/testGo/trades"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"runtime"
)

// init is invoked before main()
func init() {

	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.TextFormatter{})

	// Output to stdout instead of the default stderr
	log.SetOutput(telegram.NewWrite())

	// Set level logs.
	log.SetLevel(log.WarnLevel)

	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.WithFields(log.Fields{
			"package":  "main",
			"function": "[func init] => godotenv.Load",
			"error":    err,
		}).Fatal("Failed to get config!")
	}
}

// main function of Collector Project
func main() {
	log.Info("Start Collector")
	log.Info("Number of usage CPU : ", runtime.NumCPU())
	runtime.GOMAXPROCS(runtime.NumCPU())
	// config init
	conf := config.New()
	log.Info("Config has been successfully initialized!")

	// call Services
	go trades.StartWsTradeService(conf)
	go klines.StartAnalyticKlinesByTradesService(conf)
	go klines.StartAnalyticKlinesService(conf)
	go monitoring.StartMonitoring(conf)
	var input string
	fmt.Scanln(&input)
}
