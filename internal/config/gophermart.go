package config

import (
	"flag"
	"os"
	"time"
)

type GophermartConfig struct {
	ListenAddr           string
	DriverType           string
	LogLevel             string
	StoragePath          string
	AccrualSystemAddress string
	SecretKey            string
	TickerTime           time.Duration
}

func NewGophermartConfig() GophermartConfig {
	return loadGophermartFromFlags()
}

// Первым делом парсим из окружения
func loadGophermartFromEnv() GophermartConfig {
	var c GophermartConfig
	c.ListenAddr = os.Getenv("RUN_ADDRESS")
	c.StoragePath = os.Getenv("DATABASE_URI")
	c.AccrualSystemAddress = os.Getenv("ACCRUAL_SYSTEM_ADDRESS")
	c.SecretKey = os.Getenv("SECRET_KEY")
	return c
}

// Затем парсим из флагов. Если есть из флага то заменяем
func loadGophermartFromFlags() GophermartConfig {
	confFromEnv := loadGophermartFromEnv()
	var c GophermartConfig
	d := flag.String("d", "", "Path to store")
	l := flag.String("l", "debug", "Logger Level")
	a := flag.String("a", "", "Listen address with port")
	r := flag.String("r", "", "Accrual system address")
	k := flag.String("k", "", "Secret key for JWT")
	t := flag.Duration("t", 10*time.Second, "Ticker time")
	flag.Parse()

	c.StoragePath = ifEmpty(*d, confFromEnv.StoragePath)
	c.ListenAddr = ifEmpty(*a, confFromEnv.ListenAddr)
	c.AccrualSystemAddress = ifEmpty(*r, confFromEnv.AccrualSystemAddress)
	c.DriverType = parseDriverType(c.StoragePath)
	c.LogLevel = *l
	c.SecretKey = *k
	c.TickerTime = *t
	return c
}
