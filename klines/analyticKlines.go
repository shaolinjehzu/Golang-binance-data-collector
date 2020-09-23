package klines

import (
	"fmt"
	"github.com/shaolinjehzu/testGo/config"
	"github.com/tarantool/go-tarantool"
	log "github.com/sirupsen/logrus"
	"strconv"
	"time"
)

func FeaturesAnalyticKlines(symbol string, config *config.Config, conn *tarantool.Connection, index int){
	count, err := strconv.ParseInt(config.Data.Count_klines[index], 0, 64)
	if err != nil {

	}
	interval, err := strconv.ParseInt(config.Data.Klines[index], 0, 64)
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
	lol := dataKline[count -1].minTime

	resp, err := conn.Select(config.Data.Features_table+symbol+"_ANALYTIC_KLINES_"+config.Data.Period_klines[index-1], "primary", 0, 1000000, tarantool.IterGe, []interface{}{lol})
	if err != nil {
		log.WithFields(log.Fields{
			"package":  "klines",
			"function": "[func FeaturesAnalyticKlines] => conn.Select",
			"error":    err,
		}).Error("Failed to select klines")
	}
	log.WithFields(log.Fields{
		"package":  "klines",
		"function": "[func FeaturesAnalyticKlines] => conn.Select",
		"data":    resp.String(),
	}).Info("Selecting klines successfully")
	for _, elem := range resp.Tuples() {
		i := 0
		tt  := int64(elem[0].(uint64))
		var k int64 = 0
		for k = 0; k < count; k++ {
			if tt >= dataKline[k].minTime && tt <= dataKline[k].maxTime {
				closeP := elem[2].(float64)

				maxPrice := elem[5].(float64)

				minPrice := elem[6].(float64)

				smin := elem[3].(float64)
				smax := elem[4].(float64)

				v := elem[7].(float64)
				q := elem[8].(float64)
				qs := elem[9].(float64)
				qb := elem[10].(float64)
				n := int32(elem[11].(uint64))
				ns := int32(elem[12].(uint64))
				nb := int32(elem[13].(uint64))
				vt := elem[14].(float64)
				vm := elem[15].(float64)

				if dataKline[k].isFirst {
					klines[k].C = closeP // close price
					klines[k].L = minPrice // min price
					klines[k].H = maxPrice
					klines[k].SMIN = smin
					klines[k].SMAX = smax
					dataKline[k].isFirst = false
				}

				if maxPrice > klines[k].H {
					klines[k].H = maxPrice // max price
				}

				// get min price
				if minPrice < klines[k].L {
					klines[k].L = minPrice // min price
				}

				// max spread
				if (smax > klines[k].SMAX){
					klines[k].SMAX = smax
				}

				// min spread
				if (smin > klines[k].SMIN){
					klines[k].SMIN = smin
				}

				klines[k].V += v
				klines[k].Q += q
				klines[k].QS += qs
				klines[k].QB += qb
				klines[k].N += int32(n)
				klines[k].NS += int32(ns)
				klines[k].NB += int32(nb)
				klines[k].Vt += vt
				klines[k].Vm += vm
				dataKline[k].lastKey = i
			}
		}
		i++
	}

	for k = 0; k < count; k++ {
		resp, err = conn.Replace(config.Data.Features_table+symbol+"_ANALYTIC_KLINES_"+config.Data.Period_klines[index], []interface{}{
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
			resp, err := conn.Insert(config.Data.Features_table+symbol+"_ANALYTIC_KLINES_"+config.Data.Period_klines[index], []interface{}{
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
					"function": "[func FeaturesAnalyticKlines] => conn.Insert",
					"error":    err,
				}).Error("Failed to insert klines")
			}
			log.WithFields(log.Fields{
				"package":  "klines",
				"function": "[func FeaturesAnalyticKlines] => conn.Insert",
				"data":     resp.String(),
			}).Info("Inserting klines successfully")
		}
		log.WithFields(log.Fields{
			"package":  "klines",
			"function": "[func FeaturesAnalyticKlines] => conn.Replace",
			"data":     resp.String(),
		}).Info("Replace klines successfully")


	}
}

func SpotsAnalyticKlines(symbol string, config *config.Config, conn *tarantool.Connection, index int){
	count, err := strconv.ParseInt(config.Data.Count_klines[index], 0, 64)
	if err != nil {

	}
	interval, err := strconv.ParseInt(config.Data.Klines[index], 0, 64)
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
	lol := dataKline[count -1].minTime

	resp, err := conn.Select(config.Data.Spots_table+symbol+"_ANALYTIC_KLINES_"+config.Data.Period_klines[index-1], "primary", 0, 1000000, tarantool.IterGe, []interface{}{lol})
	if err != nil {
		log.WithFields(log.Fields{
			"package":  "klines",
			"function": "[func SpotsAnalyticKlines] => conn.Select",
			"error":    err,
		}).Error("Failed to select klines")
	}
	log.WithFields(log.Fields{
		"package":  "klines",
		"function": "[func SpotsAnalyticKline] => conn.Select",
		"data":    resp.String(),
	}).Info("Selecting klines successfully")
	for _, elem := range resp.Tuples() {
		i := 0
		tt  := int64(elem[0].(uint64))
		var k int64 = 0
		for k = 0; k < count; k++ {
			if tt >= dataKline[k].minTime && tt <= dataKline[k].maxTime {
				closeP := elem[2].(float64)

				maxPrice := elem[5].(float64)

				minPrice := elem[6].(float64)

				smin := elem[3].(float64)
				smax := elem[4].(float64)

				v := elem[7].(float64)
				q := elem[8].(float64)
				qs := elem[9].(float64)
				qb := elem[10].(float64)
				n := int32(elem[11].(uint64))
				ns := int32(elem[12].(uint64))
				nb := int32(elem[13].(uint64))
				vt := elem[14].(float64)
				vm := elem[15].(float64)

				if dataKline[k].isFirst {
					klines[k].C = closeP // close price
					klines[k].L = minPrice // min price
					klines[k].H = maxPrice
					klines[k].SMIN = smin
					klines[k].SMAX = smax
					dataKline[k].isFirst = false
				}

				if maxPrice > klines[k].H {
					klines[k].H = maxPrice // max price
				}

				// get min price
				if minPrice < klines[k].L {
					klines[k].L = minPrice // min price
				}

				// max spread
				if (smax > klines[k].SMAX){
					klines[k].SMAX = smax
				}

				// min spread
				if (smin > klines[k].SMIN){
					klines[k].SMIN = smin
				}

				klines[k].V += v
				klines[k].Q += q
				klines[k].QS += qs
				klines[k].QB += qb
				klines[k].N += int32(n)
				klines[k].NS += int32(ns)
				klines[k].NB += int32(nb)
				klines[k].Vt += vt
				klines[k].Vm += vm
				dataKline[k].lastKey = i
			}
		}
		i++
	}

	for k = 0; k < count; k++ {
		resp, err = conn.Replace(config.Data.Spots_table+symbol+"_ANALYTIC_KLINES_"+config.Data.Period_klines[index], []interface{}{
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
			resp, err := conn.Insert(config.Data.Spots_table+symbol+"_ANALYTIC_KLINES_"+config.Data.Period_klines[index], []interface{}{
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
					"function": "[func SpotsAnalyticKlines] => conn.Insert",
					"error":    err,
				}).Error("Failed to insert klines")
			}
			log.WithFields(log.Fields{
				"package":  "klines",
				"function": "[func SpotsAnalyticKlines] => conn.Insert",
				"data":     resp.String(),
			}).Info("Inserting klines successfully")
		}
		log.WithFields(log.Fields{
			"package":  "klines",
			"function": "[func SpotsAnalyticKlines] => conn.Replace",
			"data":     resp.String(),
		}).Info("Replace klines successfully")



	}
}


func StartAnalyticKlinesService(config *config.Config) {

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
		"function": "[func StartAnalyticKlines]",
	}).Info("A connection was successfully established with the Tarantool")

	symbols := config.Data.Symbols
	ticker := time.NewTicker(500 * time.Millisecond)
	for {
		<-ticker.C
		for _, symbol := range symbols {
			for i := 1; i < len(config.Data.Period_klines); i++ {
				go FeaturesAnalyticKlines(symbol, config, conn, i)
				go SpotsAnalyticKlines(symbol, config, conn, i)
			}
		}
	}
	var input string
	fmt.Scanln(&input)
}