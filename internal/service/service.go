package service

type Repository interface {
	UserRepository
	OrderRepository
	BalanceRepository
}

type ServiceConfig struct {
	Attempts         uint
	AttemptsInterval uint
}
