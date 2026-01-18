package config

import (
	"github.com/caarlos0/env"
	"github.com/spf13/pflag"
)

type Config struct {
	Address              string `env:"RUN_ADDRESS"`
	DatabaseUri          string `env:"DATABASE_URI"`
	AccrualSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
	JWTSecret            string `env:"JWT_SECRET"`
	JWTExpMinutes        uint   `env:"JWT_EXP_MINUTES"`
	RetryAttempts        uint
	RetryInterval        uint
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

	if envValues.JWTSecret == "" {
		envValues.JWTSecret = flagValues.JWTSecret
	}

	if envValues.JWTExpMinutes == 0 {
		envValues.JWTExpMinutes = flagValues.JWTExpMinutes
	}

	if envValues.RetryAttempts == 0 {
		envValues.RetryAttempts = flagValues.RetryAttempts
	}

	if envValues.RetryInterval == 0 {
		envValues.RetryInterval = flagValues.RetryInterval
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
	pflag.StringVarP(&flagValues.AccrualSystemAddress, "accrualaddr", "r", "http://localhost:3000", "Accrual system host:port")
	pflag.StringVarP(&flagValues.JWTSecret, "jwtsec", "j", "DEFAULT_SECRET", "JWT Secret") // it is not right, remove default arg
	pflag.UintVarP(&flagValues.JWTExpMinutes, "jwtexp", "s", 180, "JWT lifetime in minutes")
	pflag.UintVar(&flagValues.RetryAttempts, "retryAttempt", 3, "Count of retry attempts to execute metrics operation")
	pflag.UintVar(&flagValues.RetryInterval, "retryInterval", 2, "Interval in seconds between metric operation attempts")

	pflag.Parse()

	return &flagValues
}
