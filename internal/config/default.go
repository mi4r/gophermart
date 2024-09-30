package config

import "strings"

const (
	migrDirNameGophermart string = "gophermart"
	migrDirNameAccrual    string = "accrual"
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
