package trades

import (
	"fmt"
	"github.com/shaolinjehzu/testGo/config"
	"github.com/google/uuid"
	"github.com/shaolinjehzu/go-binance"
	log "github.com/sirupsen/logrus"
	"github.com/tarantool/go-tarantool"
)

//WS feature service handler
func featureTradeHandler(event *binance.WsTradeEvent, connection *tarantool.Connection, config *config.Config){
	resp, err := connection.Insert(config.Data.Features_table + event.Symbol + "_TRADES", []interface{}{uuid.Must(uuid.NewRandom()).String(), event.Price, event.Quantity, event.TradeTime, event.IsBuyerMaker})
	if err != nil {
		log.WithFields(log.Fields{
			"package":  "trades",
			"function": "[func featureTradeHandler] => connection.Insert",
			"error":    err,
		}).Error("Failed to insert trade")
	}
	log.WithFields(log.Fields{
		"package":  "trades",
		"function": "[func featureTradeHandler] => connection.Insert",
		"data":    resp.String(),
	}).Info("Inserting trade successfully")
}

//WS spot service handler
func spotTradeHandler(event *binance.WsTradeEvent, connection *tarantool.Connection, config *config.Config){
	resp, err := connection.Insert(config.Data.Spots_table + event.Symbol + "_TRADES", []interface{}{uuid.Must(uuid.NewRandom()).String(), event.Price, event.Quantity, event.TradeTime, event.IsBuyerMaker})
	if err != nil {
		log.WithFields(log.Fields{
			"package":  "trades",
			"function": "[func spotTradeHandler] => connection.Insert",
			"error":    err,
		}).Error("Failed to insert trade")
	}
	log.WithFields(log.Fields{
		"package":  "trades",
		"function": "[func spotTradeHandler] => connection.Insert",
		"data":    resp.String(),
	}).Info("Inserting trade successfully")
}

//WS feature service
func wsFeatureTrades(symbol string, config *config.Config, conn *tarantool.Connection) {
	wsFeatureTradeHandler := func(event *binance.WsTradeEvent) {
		go featureTradeHandler(event, conn, config)
	}
	errHandler := func(err error) {
		log.WithFields(log.Fields{
			"package":  "trades",
			"function": "[func wsSpotTrades] => errHandler",
			"error":    err,
		}).Error("Failed to get trade")
	}
	doneC, _, err := binance.WsTradeFutureServe(symbol, wsFeatureTradeHandler, errHandler)
	if err != nil {
		log.WithFields(log.Fields{
			"package":  "trades",
			"function": "[func wsSpotTrades] => binance.WsTradeFutureServe",
			"error":    err,
		}).Error("Failed to get trade")
		return
	}
	<-doneC
}

//WS spot service
func wsSpotTrades(symbol string, config *config.Config, conn *tarantool.Connection) {
	wsSpotTradeHandler := func(event *binance.WsTradeEvent) {
		go spotTradeHandler(event, conn, config)
	}
	errHandler := func(err error) {
		log.WithFields(log.Fields{
			"package":  "trades",
			"function": "[func wsSpotTrades] => errHandler",
			"error":    err,
		}).Error("Failed to get trade")
	}
	doneC, _, err := binance.WsTradeServe(symbol, wsSpotTradeHandler, errHandler)
	if err != nil {
		log.WithFields(log.Fields{
			"package":  "trades",
			"function": "[func wsSpotTrades] => binance.WsTradeServe",
			"error":    err,
		}).Error("Failed to get trade")
		return
	}
	<-doneC
}

func StartWsTradeService(config *config.Config){

	//connect to Tarantool
	conn, err := tarantool.Connect(config.Tarantool.Host, config.Tarantool.Opts)
	if err != nil {
		log.WithFields(log.Fields{
			"package":  "trades",
			"function": "[func StartWsTradeService] => tarantool.Connect",
			"error":    err,
		}).Fatal("Failed to estabilish connection!")

	}
	defer conn.Close()
	//info log if sucess
	log.WithFields(log.Fields{
		"package":  "trades",
		"function": "[func StartWsTradeService]",
	}).Info("A connection was successfully established with the Tarantool")

	symbols := config.Data.Symbols
	for _, symbol := range symbols{
		go wsFeatureTrades(symbol, config, conn)
		go wsSpotTrades(symbol, config, conn)
	}

	var input string
	fmt.Scanln(&input)
}
