package depths

import (
	"encoding/json"
	"fmt"
	"github.com/shaolinjehzu/testGo/config"
	"github.com/shaolinjehzu/go-binance"
	log "github.com/sirupsen/logrus"
	"github.com/tarantool/go-tarantool"
)

type Tick struct{
	t	int64
	symbol  string
	lBid	string
	lAsk	string
	bids	string
	asks	string

}

func SpotsDepthToTarantool(tick Tick, connection *tarantool.Connection, config *config.Config){
	resp, err := connection.Insert(config.Data.Features_table + tick.symbol + "_TICKS", []interface{}{tick.t, tick.lBid, tick.lAsk, tick.bids, tick.asks})
	if err != nil {
		log.WithFields(log.Fields{
			"package":  "depths",
			"function": "[func SpotsDepthToTarantool] => conn.Insert",
			"error":    err,
		}).Error("Failed to insert depth")
	}
	log.WithFields(log.Fields{
		"package":  "depths",
		"function": "[func SpotsDepthToTarantool] => conn.Insert",
		"data":    resp.String(),
	}).Info("Inserting depth successfully")
}

func SpotsDepth(symbol string, config *config.Config, conn *tarantool.Connection) {

	wsDepthHandler := func(event *binance.WsDepthEvent) {
		var lbid string
		var lask string
		var bids []string
		var asks []string
		var resBids []byte
		var resAsks []byte
		if len(event.Bids) != 0{
			lent := len(event.Bids) - 1
			lbid = event.Bids[lent].Price
		} else{
			lbid = "0"
		}
		if len(event.Asks) != 0{
			lent := len(event.Asks) - 1
			lask = event.Asks[lent].Price
		} else{
			lask = "0"
		}

		for _, elem := range event.Bids{
			bids = append(bids, elem.Price + ":" + elem.Quantity)
		}
		for _, elem := range event.Asks{
			asks = append(asks, elem.Price + ":" + elem.Quantity)
		}
		resBids, err := json.Marshal(bids)
		if err != nil{

		}
		resAsks, err = json.Marshal(asks)
		if err != nil{

		}

		tick := Tick{
			event.Time,
			event.Symbol,
			lbid,
			lask,
			string(resBids),
			string(resAsks),
		}
		SpotsDepthToTarantool(tick, conn, config)
	}
	errHandler := func(err error) {
		log.WithFields(log.Fields{
			"package":  "depths",
			"function": "[func SpotsDepth] => errHandler",
			"error":    err,
		}).Error("Failed to get depth")
	}
	doneC, _, err := binance.WsDepth20SpotServe(symbol, wsDepthHandler, errHandler)
	if err != nil {
		log.WithFields(log.Fields{
			"package":  "depths",
			"function": "[func SpotsDepth] => binance.WsDepth20FeatureServe",
			"error":    err,
		}).Error("Failed to get depth")
		return
	}
	<-doneC
}

func FeaturesDepthToTarantool(tick Tick, connection *tarantool.Connection, config *config.Config){
	resp, err := connection.Insert(config.Data.Features_table + tick.symbol + "_TICKS", []interface{}{tick.t, tick.lBid, tick.lAsk, tick.bids, tick.asks})
	if err != nil {
		log.WithFields(log.Fields{
			"package":  "depths",
			"function": "[func FeaturesDepthToTarantool] => conn.Insert",
			"error":    err,
		}).Error("Failed to insert depth")
	}
	log.WithFields(log.Fields{
		"package":  "depths",
		"function": "[func FeaturesDepthToTarantool] => conn.Insert",
		"data":    resp.String(),
	}).Info("Inserting depth successfully")
}

func FeaturesDepth(symbol string, config *config.Config, conn *tarantool.Connection) {

	wsDepthHandler := func(event *binance.WsDepthEvent) {
		var lbid string
		var lask string
		var bids []string
		var asks []string
		var resBids []byte
		var resAsks []byte
		if len(event.Bids) != 0{
			lent := len(event.Bids) - 1
			lbid = event.Bids[lent].Price
		} else{
			lbid = "0"
		}
		if len(event.Asks) != 0{
			lent := len(event.Asks) - 1
			lask = event.Asks[lent].Price
		} else{
			lask = "0"
		}

		for _, elem := range event.Bids{
			bids = append(bids, elem.Price + ":" + elem.Quantity)
		}
		for _, elem := range event.Asks{
			asks = append(asks, elem.Price + ":" + elem.Quantity)
		}
		resBids, err := json.Marshal(bids)
		if err != nil{

		}
		resAsks, err = json.Marshal(asks)
		if err != nil{

		}

		tick := Tick{
			event.Time,
			event.Symbol,
			lbid,
			lask,
			string(resBids),
			string(resAsks),
		}
		FeaturesDepthToTarantool(tick, conn, config)
	}
	errHandler := func(err error) {
		log.WithFields(log.Fields{
			"package":  "depths",
			"function": "[func FeaturesDepth] => errHandler",
			"error":    err,
		}).Error("Failed to get depth")
	}
	doneC, _, err := binance.WsDepth20FeatureServe(symbol, wsDepthHandler, errHandler)
	if err != nil {
		log.WithFields(log.Fields{
			"package":  "depths",
			"function": "[func FeaturesDepth] => binance.WsDepth20FeatureServe",
			"error":    err,
		}).Error("Failed to get depth")
		return
	}
	<-doneC
}


func StartWsDepthService(config *config.Config) {

	//connect to Tarantool
	conn, err := tarantool.Connect(config.Tarantool.Host, config.Tarantool.Opts)
	if err != nil {
		log.WithFields(log.Fields{
			"package":  "depths",
			"function": "[func StartWsDepthService] => tarantool.Connect",
			"error":    err,
		}).Fatal("Failed to estabilish connection!")
	}
	defer conn.Close()

	log.WithFields(log.Fields{
		"package":  "depths",
		"function": "[func StartWsDepthService]",
	}).Info("A connection was successfully established with the Tarantool")

	symbols := config.Data.Symbols
	for _, symbol := range symbols {
		go FeaturesDepth(symbol, config, conn)
		go SpotsDepth(symbol, config, conn)
	}
	
	var input string
	fmt.Scanln(&input)
}
