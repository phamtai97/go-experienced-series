package config

type DatabaseConfig struct {
	Driver                  string
	Url                     string
	ConnMaxLifetimeInMinute int
	MaxOpenConns            int
	MaxIdleConns            int
}

type HttpServerConfig struct {
	Port uint
}
