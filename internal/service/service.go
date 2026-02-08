package service

type Repository interface {
	UserRepository
	OrderRepository
	BalanceRepository
}
