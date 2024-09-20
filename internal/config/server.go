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

type Gophermart struct {
	ListenAddr           string
	DriverType           string
	LogLevel             string
	StoragePath          string
	AccrualSystemAddress string
	SecretKey            string
}

func NewGophermart() Gophermart {
	return loadGophermartFromFlags()
}

// Первым делом парсим из окружения
func loadGophermartFromEnv() Gophermart {
	var c Gophermart
	c.ListenAddr = os.Getenv("RUN_ADDRESS")
	c.StoragePath = os.Getenv("DATABASE_URI")
	c.AccrualSystemAddress = os.Getenv("ACCRUAL_SYSTEM_ADDRESS")
	c.SecretKey = os.Getenv("SECRET_KEY")
	return c
}

// Затем парсим из флагов. Если есть из флага то заменяем
func loadGophermartFromFlags() Gophermart {
	confFromEnv := loadGophermartFromEnv()
	var c Gophermart
	d := flag.String("d", "", "Path to store")
	l := flag.String("l", "debug", "Logger Level")
	a := flag.String("a", "", "Listen address with port")
	r := flag.String("r", "", "Accrual system address")
	k := flag.String("k", "", "Secret key for JWT")
	flag.Parse()

	c.StoragePath = ifEmpty(*d, confFromEnv.StoragePath)
	c.ListenAddr = ifEmpty(*a, confFromEnv.ListenAddr)
	c.AccrualSystemAddress = ifEmpty(*r, confFromEnv.AccrualSystemAddress)
	c.DriverType = parseDriverType(c.StoragePath)
	c.LogLevel = *l
	c.SecretKey = *k

	return c
}
