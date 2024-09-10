package config

import (
	"flag"
	"os"
	"strings"
)

const (
	DriverPostgres = "postgres"
)

// file://pathToStorage -> file
func parseDriverType(path string) string {
	if path == "" {
		return ""
	}
	return strings.Split(path, ":")[0]
}

func tryLoadFromEnv(key, fromFlags string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		return fromFlags
	} else {
		return value
	}
}

type ServerConfig struct {
	ListenAddr           string
	DriverType           string
	LogLevel             string
	StoragePath          string
	AccrualSystemAddress string
}

func NewServerConfig() ServerConfig {
	return loadServerConfigFromEnv()
}

func loadServerConfigFromFlags() ServerConfig {
	var config ServerConfig
	d := flag.String("d", "", "Path to store")
	l := flag.String("l", "debug", "Logger Level")
	a := flag.String("a", "localhost:8080", "Listen address with port")
	r := flag.String("r", "localhost:3000", "Accrual system address")
	flag.Parse()

	config.ListenAddr = *a
	config.StoragePath = *d
	config.DriverType = parseDriverType(config.StoragePath)
	config.AccrualSystemAddress = *r
	config.LogLevel = *l

	return config
}

func loadServerConfigFromEnv() ServerConfig {
	fromFlags := loadServerConfigFromFlags()
	fromFlags.ListenAddr = tryLoadFromEnv("RUN_ADDRESS", fromFlags.ListenAddr)
	fromFlags.StoragePath = tryLoadFromEnv("DATABASE_URI", fromFlags.StoragePath)
	fromFlags.AccrualSystemAddress = tryLoadFromEnv("ACCRUAL_SYSTEM_ADDRESS", fromFlags.AccrualSystemAddress)
	return fromFlags
}
