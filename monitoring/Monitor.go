package monitoring

import (
	"fmt"
	"github.com/shaolinjehzu/testGo/config"
	log "github.com/sirupsen/logrus"
	"github.com/tarantool/go-tarantool"
	"strconv"
	"time"
)

func CheckFeaturesTrades(symbol string, config *config.Config, conn *tarantool.Connection) {
	interval, err := strconv.ParseInt(config.Data.Klines[4], 0, 64)
	if err != nil {

	}
	unixTime := time.Now().Unix() * 1000
	balanceTime := unixTime % interval
	minTime := unixTime - balanceTime - (interval * 1)
	resp, err := conn.Select(config.Data.Features_table+symbol + "_TRADES", "secondary", 0, 1000000, tarantool.IterGe, []interface{}{(minTime)})
	if err != nil {
		log.WithFields(log.Fields{
			"package":  "monitoring",
			"function": "[func CheckFeaturesTrades] => conn.Select",
			"error":    err,
		}).Error("Failed to select trades")
	}
	if len(resp.Tuples()) == 0{
		log.Error("No Last Minute Trades for : " + symbol)
	}
}

func CheckSpotsTrades(symbol string, config *config.Config, conn *tarantool.Connection) {
	interval, err := strconv.ParseInt(config.Data.Klines[4], 0, 64)
	if err != nil {

	}
	unixTime := time.Now().Unix() * 1000
	balanceTime := unixTime % interval
	minTime := unixTime - balanceTime - (interval * 1)
	resp, err := conn.Select(config.Data.Spots_table+symbol + "_TRADES", "secondary", 0, 1000000, tarantool.IterGe, []interface{}{(minTime)})
	if err != nil {
		log.WithFields(log.Fields{
			"package":  "monitoring",
			"function": "[func CheckSpotsTrades] => conn.Select",
			"error":    err,
		}).Error("Failed to select trades")
	}
	if len(resp.Tuples()) == 0{
		log.Error("No Last Minute Trades for : " + symbol)
	}
}

func StartMonitoring(config *config.Config){
	//connect to Tarantool
	conn, err := tarantool.Connect(config.Tarantool.Host, config.Tarantool.Opts)
	if err != nil {
		log.WithFields(log.Fields{
			"package":  "monitoring",
			"function": "[func CheckFeaturesTrades] => tarantool.Connect",
			"error":    err,
		}).Fatal("Failed to estabilish connection!")
	}
	defer conn.Close()

	log.WithFields(log.Fields{
		"package":  "klines",
		"function": "[func StartAnalyticKlines]",
	}).Info("A connection was successfully established with the Tarantool")

	symbols := config.Data.Symbols
	ticker := time.NewTicker(59000 * time.Millisecond)
	for {
		<-ticker.C
		for _, symbol := range symbols {
			go CheckFeaturesTrades(symbol, config, conn)
		}
	}
	var input string
	fmt.Scanln(&input)
}
