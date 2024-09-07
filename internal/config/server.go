package config

type ServerConfig struct {
	ListenAddr string
	DriverType string
	LogLevel   string
}

func NewServerConfig() ServerConfig {
	return ServerConfig{
		ListenAddr: ":8080",
		DriverType: "memory",
		LogLevel:   "debug",
	}
}
