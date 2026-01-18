package config

import (
	"github.com/caarlos0/env"
	"github.com/spf13/pflag"
)

type Config struct {
	Address              string `env:"RUN_ADDRESS"`
	DatabaseUri          string `env:"DATABASE_URI"`
	AccrualSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
}

func InitConfig() *Config {
	envValues := initEnv()
	flagValues := initFlags()

	if envValues.Address == "" {
		envValues.Address = flagValues.Address
	}

	if envValues.DatabaseUri == "" {
		envValues.DatabaseUri = flagValues.DatabaseUri
	}

	if envValues.AccrualSystemAddress == "" {
		envValues.AccrualSystemAddress = flagValues.AccrualSystemAddress
	}

	return envValues
}

func initEnv() *Config {
	agentConfig := Config{}

	env.Parse(&agentConfig)

	return &agentConfig
}

func initFlags() *Config {
	flagValues := Config{}

	pflag.StringVarP(&flagValues.Address, "addr", "a", "localhost:8080", "Address host:port")
	pflag.StringVarP(&flagValues.DatabaseUri, "dbaddr", "d", "postgresql://localhost/postgres", "PG URI")
	pflag.StringVarP(&flagValues.AccrualSystemAddress, "accuraladdr", "r", "http://localhost:3000", "Accural system host:port")

	pflag.Parse()

	return &flagValues
}
