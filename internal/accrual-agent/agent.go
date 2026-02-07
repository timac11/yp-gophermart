package accrualagent

import (
	"context"
	"sync"
	"time"

	"sync/atomic"

	"github.com/timac11/yp-gophermart/internal/logger"
	"github.com/timac11/yp-gophermart/internal/model"
	"go.uber.org/zap"
)

const (
	timeoutGetOrdersDB = 5
)

type Repository interface {
	GetProcessingAccruals(ctx context.Context, limit uint32, offset uint32) ([]*model.Accrual, error)
	UpdateAccrualStatus(ctx context.Context, accrual model.Accrual) error
}

type WorkerTimeoutsConfig struct {
	ClientTimeout     time.Duration // http client timeout per request
	RetryTimeout      time.Duration // timeout between two retries
	ProcessingTimeout time.Duration // timeout for whole polling task
}

type AgentLimits struct {
	WorkersCount uint
	OrdersBuffer uint32
}

type Agent struct {
	repository            Repository
	accrualURL            string
	chOrdersForProcessing chan string
	chOrdersResult        chan model.Accrual
	workersTimeoutConfig  WorkerTimeoutsConfig
	limits                AgentLimits
	alreadyProcessedCount atomic.Uint32
}

func NewAgent(
	repository Repository,
	accrualURL string,
	workersTimeoutConfig WorkerTimeoutsConfig,
	limits AgentLimits,
) *Agent {
	return &Agent{
		repository:            repository,
		accrualURL:            accrualURL,
		workersTimeoutConfig:  workersTimeoutConfig,
		limits:                limits,
		chOrdersForProcessing: make(chan string, limits.OrdersBuffer),
		chOrdersResult:        make(chan model.Accrual, limits.OrdersBuffer),
	}
}

func (agent *Agent) Run(ctx context.Context) {
	go agent.runGetActualOrders(ctx)
	go agent.runWorkers(ctx)
	go agent.runUpdateAccrualStatuses(ctx)
}

func (agent *Agent) runGetActualOrders(ctx context.Context) {
	agent.getProcessingAccruals(ctx)

	ticker := time.NewTicker(timeoutGetOrdersDB * time.Second)

	defer close(agent.chOrdersForProcessing)

	for {
		select {
		case <-ticker.C:
			agent.getProcessingAccruals(ctx)
		case <-ctx.Done():
			return
		}
	}
}

func (agent *Agent) getProcessingAccruals(ctx context.Context) {
	log := logger.LoggerFromContext(ctx)

	// initial load
	orders, err := agent.repository.GetProcessingAccruals(ctx, agent.limits.OrdersBuffer, agent.alreadyProcessedCount.Load())
	if err != nil {
		log.Error(err.Error())
	} else {
		agent.alreadyProcessedCount.Add(uint32(len(orders)))
		for _, order := range orders {
			agent.chOrdersForProcessing <- order.Order
		}
	}
}

func (agent *Agent) runWorkers(ctx context.Context) {
	workersCount := agent.limits.WorkersCount
	var wg sync.WaitGroup

	for i := 1; i <= int(workersCount); i++ {
		wg.Add(1)

		go func(id int) {
			defer wg.Done()
			worker := NewWorker(
				agent.accrualURL,
				agent.chOrdersForProcessing,
				agent.chOrdersResult,
				agent.workersTimeoutConfig,
			)
			worker.Run(ctx)
		}(i)
	}

	defer close(agent.chOrdersResult)
	wg.Wait()
}

func (agent *Agent) runUpdateAccrualStatuses(ctx context.Context) {
	log := logger.LoggerFromContext(ctx)

	for orderUpdate := range agent.chOrdersResult {
		log.Info("Update order:",
			zap.String("orderId", orderUpdate.Order),
			zap.String("status", string(orderUpdate.Status)),
		)
		err := agent.repository.UpdateAccrualStatus(ctx, orderUpdate)
		if err != nil {
			log.Error(err.Error())
		}
	}
}
