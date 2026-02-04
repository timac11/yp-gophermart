package accrualagent

import (
	"context"
	"time"

	"github.com/timac11/yp-gophermart/internal/logger"
	"github.com/timac11/yp-gophermart/internal/model"
)

const (
	timeoutGetOrdersDB = 5
)

type Repository interface {
	GetProcessingOrders(ctx context.Context, lastUpdateTime *time.Time) ([]*model.OrderModel, error)
	UpdateAccrualStatus(ctx context.Context) error
}

type WorkerTimeoutsConfig struct {
	clientTimeout     time.Duration // http client timeout per request
	retryTimeout      time.Duration // timeout between two retries
	processingTimeout time.Duration // timeout for whole polling task
}

type AgentLimits struct {
	workersCount int
	ordersBuffer int
}

type Agent struct {
	repository            Repository
	accrualURL            string
	lastOrdersFetchFromDb *time.Time
	chOrdersForProcessing chan string
	chOrdersResult        chan model.Accrual
	workersTimeoutConfig  WorkerTimeoutsConfig
	limits                AgentLimits
}

func NewAgent(
	repository Repository,
	accrualURL string,
	workersTimeoutConfig WorkerTimeoutsConfig,
	limits AgentLimits,
) *Agent {
	return &Agent{
		repository:           repository,
		accrualURL:           accrualURL,
		workersTimeoutConfig: workersTimeoutConfig,
		limits:               limits,
	}
}

func (agent *Agent) Start(ctx context.Context) {
	go agent.runGetActualOrders(ctx)
	go agent.runWorkers(ctx)
	go agent.runUpdateAccrualStatuses(ctx)
}

func (agent *Agent) runGetActualOrders(ctx context.Context) {
	ticker := time.NewTicker(timeoutGetOrdersDB * time.Second)
	log := logger.LoggerFromContext(ctx)

	for {
		select {
		case <-ticker.C:
			orders, err := agent.repository.GetOrders(ctx, agent.lastOrdersFetchFromDb)
			if err != nil {
				log.Error(err.Error())
			} else {
				now := time.Now()
				agent.lastOrdersFetchFromDb = &now
				for _, order := range orders {
					agent.chOrdersForProcessing <- order.OrderNum
				}
			}
		case <-ctx.Done():
			return
		}
	}
}

func (agent *Agent) runWorkers(ctx context.Context) {}

func (agent *Agent) runUpdateAccrualStatuses(ctx context.Context) {}
