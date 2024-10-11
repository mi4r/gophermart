package config

import (
	"flag"
	"os"
)

type AccrualConfig struct {
	ListenAddr  string
	DriverType  string
	LogLevel    string
	StoragePath string
	RateLimit   int
}

func NewAccrualConfig() AccrualConfig {
	return loadAccSysConfigFromFlags()
}

// Первым делом парсим из окружения
func loadAccrualConfigFromEnv() AccrualConfig {
	var c AccrualConfig
	c.ListenAddr = os.Getenv("RUN_ADDRESS")
	c.StoragePath = os.Getenv("DATABASE_URI")
	return c
}

// Затем парсим из флагов. Если есть из флага то заменяем
func loadAccSysConfigFromFlags() AccrualConfig {
	confFromEnv := loadAccrualConfigFromEnv()
	var c AccrualConfig
	d := flag.String("d", "", "Path to store")
	l := flag.String("l", "debug", "Logger Level")
	a := flag.String("a", "", "Listen address with port")
	flag.Parse()

	c.StoragePath = ifEmpty(*d, confFromEnv.StoragePath)
	c.ListenAddr = ifEmpty(*a, confFromEnv.ListenAddr)

	c.DriverType = parseDriverType(c.StoragePath)
	c.LogLevel = *l

	// 5 запросов в секунду
	c.RateLimit = 5

	return c
}
