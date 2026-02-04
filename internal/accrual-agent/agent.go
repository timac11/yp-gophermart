package accrualagent

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/timac11/yp-gophermart/internal/logger"
	"github.com/timac11/yp-gophermart/internal/model"
	"go.uber.org/zap"
)

const (
	timeoutClient       = 5
	maxWorkers          = 3
	bufSizeOrdersRecord = 3
	limitQuery          = 10
	timeoutLoadOrdersDB = 3
	timeoutGetOrdersDB  = 5
)

type Repository interface {
	GetOrders(ctx context.Context, lastUpdateTime *time.Time) ([]*model.OrderModel, error)
	UpdateOrders(ctx context.Context) error
}

type Agent struct {
	repository Repository
	accrualURL string
	// DB polling settings
	lastOrdersFetchFromDb *time.Time
	ordersFetchLimit      int
	chOrdersForProcessing chan string
	// worker settings

}

func NewAgent(repository Repository, accrualURL string) *Agent {
	return &Agent{
		repository: repository,
		accrualURL: accrualURL,
		client:     &http.Client{Timeout: time.Second * timeoutClient},
	}
}

func (agent *Agent) Start(ctx context.Context) {

}

func (agent *Agent) runGetActualOrders(ctx context.Context) {
	ticker := time.NewTicker(timeoutGetOrdersDB * time.Second)
	log := logger.LoggerFromContext(ctx)

	for {
		select {
		case <-ticker.C:
			orders, err := agent.repository.GetOrders(ctx, nil)
			if err != nil {
				log.Error(err.Error())
			} else {
				for _, order := range orders {
					agent.chOrdersForProcessing <- order
				}
			}
		case <-ctx.Done():
			return
		}
	}
}

func (agent *Agent) runGetOrdersInfoFromAccrualService(ctx context.Context) {

}

func (agent *Agent) runGetOrderInfoWorker(ctx context.Context, worker int) {
	log := logger.LoggerFromContext(ctx)

	for orderInfo := range agent.chOrdersForProcessing {
		log.Info("start getting info", zap.String("orderNum", orderInfo.OrderNum))
		url := fmt.Sprintf("%s%s%d", agent.accrualURL, "/api/orders/", orderInfo.OrderNum)
		resp, err := agent.client.Get(url)

		if err != nil {
			log.Error("Get order status error: " + err.Error())
		} else {
			// todo: parse response to structure
		}
	}
}

func (agent *Agent) runUpdateOrdersStatus(ctx context.Context) {

}
