package config

import (
	"github.com/tarantool/go-tarantool"
	"os"
	"strconv"
	"strings"
)

type TarantoolConfig struct {
	Host     string
	Opts     tarantool.Opts
}


type DataConfig struct{
	Symbols 	   			 []string
	Klines  	   			 []string
	Count_klines   			 []string
	Period_klines  			 []string
	Features_table	   		 string
	Spots_table 	   		 string
	Data_types				 []string
}

type Config struct {
	Tarantool TarantoolConfig
	Data DataConfig
}

// New returns a new Config struct
func New() *Config {
	return &Config{
		Tarantool: TarantoolConfig{
			Host:     getEnv("TARANTOOL_HOST", "127.0.0.1:3301"),
			Opts:     tarantool.Opts{
							User: getEnv("TARANTOOL_USERNAME", "guest"),
							Pass: getEnv("TARANTOOL_PASSWORD", ""),
						},
		},
		Data: DataConfig{
			Symbols: getEnvAsSlice("SYMBOLS", []string{"BTCUSDT"}, ","),
			Klines:  getEnvAsSlice("KLINES", []string{"BTCUSDT"}, ","),
			Period_klines: getEnvAsSlice("PERIOD_KLINES", []string{"100ms"}, ","),
			Count_klines:  getEnvAsSlice("COUNT_KLINES", []string{"BTCUSDT"}, ","),
			Features_table:   getEnv("BINANCE_FEATURES_UNIVERSE_TABLE", "BINANCE_FEATURES_"),
			Spots_table: getEnv("BINANCE_SPOTS_UNIVERSE_TABLE", "BINANCE_SPOT_"),
			Data_types: getEnvAsSlice("DATA_TYPES", []string{"FEATURES","SPOT"}, ","),
		},
	}
}

// Simple helper function to read an environment or return a default value
func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}

// Simple helper function to read an environment variable into integer or return a default value
func getEnvAsInt(name string, defaultVal int) int {
	valueStr := getEnv(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}

	return defaultVal
}

// Helper to read an environment variable into a bool or return default value
func getEnvAsBool(name string, defaultVal bool) bool {
	valStr := getEnv(name, "")
	if val, err := strconv.ParseBool(valStr); err == nil {
		return val
	}

	return defaultVal
}

// Helper to read an environment variable into a string slice or return default value
func getEnvAsSlice(name string, defaultVal []string, sep string) []string {
	valStr := getEnv(name, "")

	if valStr == "" {
		return defaultVal
	}

	val := strings.Split(valStr, sep)

	return val
}