package config

import (
	"flag"
	"os"
)

type AccSysConfig struct {
	ListenAddr  string
	DriverType  string
	LogLevel    string
	StoragePath string
	MigrDirName string
}

func NewAccSysConfig() AccSysConfig {
	return loadAccSysConfigFromFlags()
}

// Первым делом парсим из окружения
func loadAccSysConfigFromEnv() AccSysConfig {
	var c AccSysConfig
	c.ListenAddr = os.Getenv("RUN_ADDRESS")
	c.StoragePath = os.Getenv("DATABASE_URI")
	return c
}

// Затем парсим из флагов. Если есть из флага то заменяем
func loadAccSysConfigFromFlags() AccSysConfig {
	confFromEnv := loadAccSysConfigFromEnv()
	var c AccSysConfig
	d := flag.String("d", "", "Path to store")
	l := flag.String("l", "debug", "Logger Level")
	a := flag.String("a", "", "Listen address with port")
	flag.Parse()

	c.StoragePath = ifEmpty(*d, confFromEnv.StoragePath)
	c.ListenAddr = ifEmpty(*a, confFromEnv.ListenAddr)

	c.DriverType = parseDriverType(c.StoragePath)
	c.LogLevel = *l

	c.MigrDirName = migrDirNameAccrual

	return c
}
