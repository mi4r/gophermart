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

func ifEmpty(fromFlag, fromEnv string) string {
	// Если пусто во флаге, то вернем из окружения
	if fromFlag == "" {
		return fromEnv
	}
	return fromFlag
}

type ServerConfig struct {
	ListenAddr           string
	DriverType           string
	LogLevel             string
	StoragePath          string
	AccrualSystemAddress string
}

func NewServerConfig() ServerConfig {
	return loadServerConfigFromFlags()
}

// Первым делом парсим из окружения
func loadServerConfigFromEnv() ServerConfig {
	var c ServerConfig
	c.ListenAddr = os.Getenv("RUN_ADDRESS")
	c.StoragePath = os.Getenv("DATABASE_URI")
	c.AccrualSystemAddress = os.Getenv("ACCRUAL_SYSTEM_ADDRESS")
	return c
}

// Затем парсим из флагов. Если есть из флага то заменяем
func loadServerConfigFromFlags() ServerConfig {
	confFromEnv := loadServerConfigFromEnv()
	var c ServerConfig
	d := flag.String("d", "", "Path to store")
	l := flag.String("l", "debug", "Logger Level")
	a := flag.String("a", "", "Listen address with port")
	r := flag.String("r", "", "Accrual system address")
	flag.Parse()

	c.StoragePath = ifEmpty(*d, confFromEnv.StoragePath)
	c.ListenAddr = ifEmpty(*a, confFromEnv.ListenAddr)
	c.AccrualSystemAddress = ifEmpty(*r, confFromEnv.AccrualSystemAddress)
	c.DriverType = parseDriverType(c.StoragePath)
	c.LogLevel = *l

	return c
}
