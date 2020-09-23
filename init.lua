#!/usr/bin/env tarantool

box.cfg {
	listen = 3301,
	net_msg_max=15000,
	readahead=240000,
	memtx_memory=512 * 1024 *1024
}

box.schema.user.passwd('pass')
--box.schema.user.grant('guest','read,write,execute,create,drop','universe')

function function_trades()
	local types={"FEATURES", "SPOT"}
	local arr={"BTCUSDT", "ETHUSDT","BCHUSDT","XRPUSDT","EOSUSDT","LTCUSDT","TRXUSDT","ETCUSDT","LINKUSDT","XLMUSDT","ADAUSDT","XMRUSDT","DASHUSDT","ZECUSDT","XTZUSDT","BNBUSDT","ATOMUSDT","ONTUSDT","IOTAUSDT","BATUSDT","VETUSDT","NEOUSDT","QTUMUSDT","IOSTUSDT","THETAUSDT","ALGOUSDT","ZILUSDT","KNCUSDT","ZRXUSDT","COMPUSDT","OMGUSDT","DOGEUSDT","SXPUSDT","LENDUSDT"}
	for i = 1, #arr do
		for k = 1, #types do
			s=box.schema.space.create("BINANCE_".. types[k] .."_" .. arr[i] .. "_TRADES")
			s:format({
				{name='id',type='string'},
				{name='price',type='string'},
				{name='quantity',type='string'},
				{name='time',type='unsigned'},
				{name='isbuyermaker',type='boolean'}
			})
			s:create_index('primary',{type='HASH',parts={'id'}})
			s:create_index('secondary',{unique=false, if_not_exists=true, type='TREE',parts={'time'}})
		end
	end
end

function function_klines()
	local types={"FEATURES", "SPOT"}
	local periods = {"100MS","1S","10S","30S","1M","3M","5M","15M","30M","1H","2H","4H","8H","1D","1W","4W"}
	local arr={"BTCUSDT", "ETHUSDT","BCHUSDT","XRPUSDT","EOSUSDT","LTCUSDT","TRXUSDT","ETCUSDT","LINKUSDT","XLMUSDT","ADAUSDT","XMRUSDT","DASHUSDT","ZECUSDT","XTZUSDT","BNBUSDT","ATOMUSDT","ONTUSDT","IOTAUSDT","BATUSDT","VETUSDT","NEOUSDT","QTUMUSDT","IOSTUSDT","THETAUSDT","ALGOUSDT","ZILUSDT","KNCUSDT","ZRXUSDT","COMPUSDT","OMGUSDT","DOGEUSDT","SXPUSDT","LENDUSDT"}
	for i = 1, #arr do
		for k = 1, #periods do
			for l = 1, #types do
				s=box.schema.space.create("BINANCE_".. types[l] .."_" .. arr[i] .."_ANALYTIC_KLINES_" .. periods[k])
				s:format({
					{name='t',type='unsigned'},
					{name='o',type='double'},
					{name='c',type='double'},
					{name='smin',type='double'},
					{name='smax',type='double'},
					{name='h',type='double'},
					{name='l',type='double'},
					{name='v',type='double'},
					{name='q',type='double'},
					{name='qs',type='double'},
					{name='qb',type='double'},
					{name='n',type='unsigned'},
					{name='ns',type='unsigned'},
					{name='nb',type='unsigned'},
					{name='Vt',type='double'},
					{name='Vm',type='double'},
				})
				s:create_index('primary',{type='TREE',parts={'t'}})
			end
		end
	end
end

function function_ticks()
	local types={"FEATURES", "SPOT"}
	local arr={"BTCUSDT", "ETHUSDT","BCHUSDT","XRPUSDT","EOSUSDT","LTCUSDT","TRXUSDT","ETCUSDT","LINKUSDT","XLMUSDT","ADAUSDT","XMRUSDT","DASHUSDT","ZECUSDT","XTZUSDT","BNBUSDT","ATOMUSDT","ONTUSDT","IOTAUSDT","BATUSDT","VETUSDT","NEOUSDT","QTUMUSDT","IOSTUSDT","THETAUSDT","ALGOUSDT","ZILUSDT","KNCUSDT","ZRXUSDT","COMPUSDT","OMGUSDT","DOGEUSDT","SXPUSDT","LENDUSDT"}
	for i = 1, #arr do
		for l = 1, #types do
			s=box.schema.space.create("BINANCE_".. types[l] .. "_" .. arr[i] .."_TICKS")
			s:format({
				{name='t',type='unsigned'},
				{name='lbid',type='string'},
				{name='lask',type='string'},
				{name='bids',type='string'},
				{name='asks',type='string'},
			})
			s:create_index('primary',{type='TREE',parts={'t'}})
		end
	end
end


local function bootstrap()
	function_trades()
	function_klines()
	function_ticks()
end

-- for first run create a space and add set up grants
box.once('GOLANG', bootstrap)

