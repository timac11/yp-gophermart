package config

import (
	"time"

	"github.com/caarlos0/env"
	"github.com/spf13/pflag"
)

type Config struct {
	Address              string `env:"RUN_ADDRESS"`
	DatabaseURI          string `env:"DATABASE_URI"`
	AccrualSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
	JWTSecret            string `env:"JWT_SECRET"`
	JWTExpMinutes        uint   `env:"JWT_EXP_MINUTES"`
	ClientTimeout        time.Duration
	RetryTimeout         time.Duration
	ProcessingTimeout    time.Duration
	WorkersCount         uint
	OrdersBuffer         uint32
}

func InitConfig() *Config {
	envValues := initEnv()
	flagValues := initFlags()

	if envValues.Address == "" {
		envValues.Address = flagValues.Address
	}

	if envValues.DatabaseURI == "" {
		envValues.DatabaseURI = flagValues.DatabaseURI
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

	envValues.ClientTimeout = time.Second * 2
	envValues.RetryTimeout = time.Second * 1
	envValues.ProcessingTimeout = time.Second * 3

	envValues.WorkersCount = 3
	envValues.OrdersBuffer = 10

	return envValues
}

func initEnv() *Config {
	agentConfig := Config{}

	env.Parse(&agentConfig)

	return &agentConfig
}

func initFlags() *Config {
	flagValues := Config{}

	pflag.StringVarP(&flagValues.Address, "addr", "a", "localhost:3000", "Address host:port")
	pflag.StringVarP(&flagValues.DatabaseURI, "dbaddr", "d", "postgresql://localhost/postgres", "PG URI")
	pflag.StringVarP(&flagValues.AccrualSystemAddress, "accrualaddr", "r", "http://localhost:8080", "Accrual system host:port")
	pflag.StringVarP(&flagValues.JWTSecret, "jwtsec", "j", "DEFAULT_SECRET", "JWT Secret") // it is not right, remove default arg
	pflag.UintVarP(&flagValues.JWTExpMinutes, "jwtexp", "s", 180, "JWT lifetime in minutes")
	pflag.Parse()

	return &flagValues
}
