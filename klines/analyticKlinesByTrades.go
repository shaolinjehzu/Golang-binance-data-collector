package klines

import (
	"fmt"
	"github.com/shaolinjehzu/testGo/config"
	log "github.com/sirupsen/logrus"
	"github.com/tarantool/go-tarantool"
	"strconv"
	"time"
)

func SpotsAnalyticKlinesByTrades(symbol string, config *config.Config, conn *tarantool.Connection) {
	count, err := strconv.ParseInt(config.Data.Count_klines[0], 0, 64)
	if err != nil {

	}
	interval, err := strconv.ParseInt(config.Data.Klines[0], 0, 64)
	if err != nil {

	}
	var k int64 = 0
	dataKline := make([]dataKline, count)
	klines := make([]Klines, count)
	unixTime := time.Now().Unix() * 1000
	balanceTime := unixTime % interval
	for k = 0; k < count; k++ {
		dataKline[k].minTime = unixTime - balanceTime - (interval * k)
		dataKline[k].maxTime = dataKline[k].minTime + interval - 1
		dataKline[k].isFirst = true
		dataKline[k].lastKey = 0

		klines[k].T = dataKline[k].minTime
		klines[k].O = 0
		klines[k].C = 0
		klines[k].H = 0
		klines[k].L = 0
		klines[k].V = 0
		klines[k].Q = 0
		klines[k].QS = 0
		klines[k].QB = 0
		klines[k].N = 0
		klines[k].N = 0
		klines[k].NB = 0
		klines[k].Vt = 0
		klines[k].Vm = 0
	}

	resp, err := conn.Select(config.Data.Spots_table+symbol+"_TRADES", "secondary", 0, 1000000, tarantool.IterGe, []interface{}{(dataKline[count-1].minTime)})
	if err != nil {
		log.WithFields(log.Fields{
			"package":  "klines",
			"function": "[func SpotsAnalyticKlinesByTrades] => conn.Select",
			"error":    err,
		}).Warn("Failed to select trades")
	}
	log.WithFields(log.Fields{
		"package":  "klines",
		"function": "[func SpotsAnalyticKlinesByTrades] => conn.Select",
		"data":     resp.String(),
	}).Info("Selecting trades successfully")
	trades := make([]trades, len(resp.Tuples()))
	if len(resp.Tuples()) > 0 {
		for _, elem := range resp.Tuples() {
			i := 0
			var tt int64 = int64(elem[3].(uint64))
			var k int64 = 0
			for k = 0; k < count; k++ {

				closeP, err := strconv.ParseFloat(elem[1].(string), 64)
				trades[i].o = closeP
				if err != nil {

				}
				volume, err := strconv.ParseFloat(elem[2].(string), 64)
				if err != nil {

				}
				if tt >= dataKline[k].minTime && tt <= dataKline[k].maxTime {

					if dataKline[k].isFirst {

						klines[k].C = closeP // close price
						klines[k].L = closeP // min price
						klines[k].H = closeP // max price
						dataKline[k].isFirst = false
					}

					if closeP > klines[k].H {
						klines[k].H = closeP // max price
					}

					// get min price
					if closeP < klines[k].L {
						klines[k].L = closeP // min price
					}

					// volume
					klines[k].V += volume

					//quote volume
					klines[k].Q += volume * closeP

					if elem[4] == "1" { // buyer
						klines[k].Vm += volume
						klines[k].QB += volume * closeP
						klines[k].NB++
					} else { // seller
						klines[k].Vt += volume
						klines[k].QS += volume * closeP
						klines[k].NS++
					}

					klines[k].N++
					dataKline[k].lastKey = i
				}
			}
			i++
		}

		for k = 0; k < count; k++ {
			klines[k].O = trades[dataKline[k].lastKey].o
			resp, err = conn.Insert(config.Data.Spots_table+symbol+"_ANALYTIC_KLINES_"+config.Data.Period_klines[0], []interface{}{
				klines[k].T,
				klines[k].O,
				klines[k].C,
				klines[k].SMIN,
				klines[k].SMAX,
				klines[k].H,
				klines[k].L,
				klines[k].V,
				klines[k].Q,
				klines[k].QS,
				klines[k].QB,
				klines[k].N,
				klines[k].NS,
				klines[k].NB,
				klines[k].Vt,
				klines[k].Vm,
			})
			if err != nil {
				resp, err := conn.Replace(config.Data.Spots_table+symbol+"_ANALYTIC_KLINES_"+config.Data.Period_klines[0], []interface{}{
					klines[k].T,
					klines[k].O,
					klines[k].C,
					klines[k].SMIN,
					klines[k].SMAX,
					klines[k].H,
					klines[k].L,
					klines[k].V,
					klines[k].Q,
					klines[k].QS,
					klines[k].QB,
					klines[k].N,
					klines[k].NS,
					klines[k].NB,
					klines[k].Vt,
					klines[k].Vm,
				})
				if err != nil {
					log.WithFields(log.Fields{
						"package":  "klines",
						"function": "[func SpotsAnalyticKlinesByTrades] => conn.Replace",
						"error":    err,
					}).Error("Failed to Replace klines")
				}
				log.WithFields(log.Fields{
					"package":  "klines",
					"function": "[func SpotsAnalyticKlinesByTrades] => conn.Replace",
					"data":     resp.String(),
				}).Info("Replacing klines successfully")
			}
			log.WithFields(log.Fields{
				"package":  "klines",
				"function": "[func SpotsAnalyticKlinesByTrades] => conn.Insert",
				"data":     resp.String(),
			}).Info("Inserting klines successfully")

		}
	}
}

func FeatureAnalyticKlinesByTrades(symbol string, config *config.Config, conn *tarantool.Connection) {
	count, err := strconv.ParseInt(config.Data.Count_klines[0], 0, 64)
	if err != nil {

	}
	interval, err := strconv.ParseInt(config.Data.Klines[0], 0, 64)
	if err != nil {

	}
	var k int64 = 0
	dataKline := make([]dataKline, count)
	klines := make([]Klines, count)
	unixTime := time.Now().Unix() * 1000
	balanceTime := unixTime % interval
	for k = 0; k < count; k++ {
		dataKline[k].minTime = unixTime - balanceTime - (interval * k)
		dataKline[k].maxTime = dataKline[k].minTime + interval - 1
		dataKline[k].isFirst = true
		dataKline[k].lastKey = 0

		klines[k].T = dataKline[k].minTime
		klines[k].O = 0
		klines[k].C = 0
		klines[k].H = 0
		klines[k].L = 0
		klines[k].V = 0
		klines[k].Q = 0
		klines[k].QS = 0
		klines[k].QB = 0
		klines[k].N = 0
		klines[k].N = 0
		klines[k].NB = 0
		klines[k].Vt = 0
		klines[k].Vm = 0
	}

	resp, err := conn.Select(config.Data.Features_table+symbol+"_TRADES", "secondary", 0, 1000000, tarantool.IterGe, []interface{}{(dataKline[count-1].minTime)})
	if err != nil {
		log.WithFields(log.Fields{
			"package":  "klines",
			"function": "[func FeatureAnalyticKlinesByTrades] => conn.Select",
			"error":    err,
		}).Error("Failed to select klines")
	}
	log.WithFields(log.Fields{
		"package":  "klines",
		"function": "[func FeatureAnalyticKlinesByTrades] => conn.Select",
		"data":     resp.String(),
	}).Info("Selecting klines successfully")
	if len(resp.Tuples()) > 0 {
		trades := make([]trades, len(resp.Tuples()))
		for _, elem := range resp.Tuples() {
			i := 0
			var tt int64 = int64(elem[3].(uint64))
			var k int64 = 0
			for k = 0; k < count; k++ {

				closeP, err := strconv.ParseFloat(elem[1].(string), 64)
				if err != nil {

				}
				trades[i].o = closeP
				volume, err := strconv.ParseFloat(elem[2].(string), 64)
				if err != nil {

				}
				if tt >= dataKline[k].minTime && tt <= dataKline[k].maxTime {

					if dataKline[k].isFirst {

						klines[k].C = closeP // close price
						klines[k].L = closeP // min price
						klines[k].H = closeP // max price
						dataKline[k].isFirst = false
					}

					if closeP > klines[k].H {
						klines[k].H = closeP // max price
					}

					// get min price
					if closeP < klines[k].L {
						klines[k].L = closeP // min price
					}

					// volume
					klines[k].V += volume

					//quote volume
					klines[k].Q += volume * closeP

					if elem[4] == "1" { // buyer
						klines[k].Vm += volume
						klines[k].QB += volume * closeP
						klines[k].NB++
					} else { // seller
						klines[k].Vt += volume
						klines[k].QS += volume * closeP
						klines[k].NS++
					}

					klines[k].N++
					dataKline[k].lastKey = i
				}
			}
			i++
		}

		for k = 0; k < count; k++ {
			klines[k].O = trades[dataKline[k].lastKey].o
			resp, err = conn.Insert(config.Data.Spots_table+symbol+"_ANALYTIC_KLINES_"+config.Data.Period_klines[0], []interface{}{
				klines[k].T,
				klines[k].O,
				klines[k].C,
				klines[k].SMIN,
				klines[k].SMAX,
				klines[k].H,
				klines[k].L,
				klines[k].V,
				klines[k].Q,
				klines[k].QS,
				klines[k].QB,
				klines[k].N,
				klines[k].NS,
				klines[k].NB,
				klines[k].Vt,
				klines[k].Vm,
			})
			if err != nil {
				resp, err := conn.Replace(config.Data.Features_table+symbol+"_ANALYTIC_KLINES_"+config.Data.Period_klines[0], []interface{}{
					klines[k].T,
					klines[k].O,
					klines[k].C,
					klines[k].SMIN,
					klines[k].SMAX,
					klines[k].H,
					klines[k].L,
					klines[k].V,
					klines[k].Q,
					klines[k].QS,
					klines[k].QB,
					klines[k].N,
					klines[k].NS,
					klines[k].NB,
					klines[k].Vt,
					klines[k].Vm,
				})
				if err != nil {
					log.WithFields(log.Fields{
						"package":  "klines",
						"function": "[func FeaturesAnalyticKlinesByTrades] => conn.Replace",
						"error":    err,
					}).Error("Failed to Replace klines")
				}
				log.WithFields(log.Fields{
					"package":  "klines",
					"function": "[func FeaturesAnalyticKlinesByTrades] => conn.Replace",
					"data":     resp.String(),
				}).Info("Replacing klines successfully")
			}
			log.WithFields(log.Fields{
				"package":  "klines",
				"function": "[func FeaturesAnalyticKlinesByTrades] => conn.Insert",
				"data":     resp.String(),
			}).Info("Inserting klines successfully")

		}
	}
}

func StartAnalyticKlinesByTradesService(config *config.Config) {
	//connect to Tarantool
	conn, err := tarantool.Connect(config.Tarantool.Host, config.Tarantool.Opts)
	if err != nil {
		log.WithFields(log.Fields{
			"package":  "klines",
			"function": "[func StartAnalyticKlines] => tarantool.Connect",
			"error":    err,
		}).Fatal("Failed to estabilish connection!")
	}
	defer conn.Close()

	log.WithFields(log.Fields{
		"package":  "klines",
		"function": "[func StartAnalyticKlinesByTrades]",
	}).Info("A connection was successfully established with the Tarantool")

	symbols := config.Data.Symbols
	ticker := time.NewTicker(500 * time.Millisecond)
	for {
		<-ticker.C
		for _, symbol := range symbols {
			go FeatureAnalyticKlinesByTrades(symbol, config, conn)
			go SpotsAnalyticKlinesByTrades(symbol, config, conn)
		}
	}
	var input string
	fmt.Scanln(&input)
}
