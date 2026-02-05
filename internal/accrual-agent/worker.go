package accrualagent

import (
	"context"
	"encoding/json"
	_errors "errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/avast/retry-go/v4"

	"github.com/timac11/yp-gophermart/internal/errors"
	"github.com/timac11/yp-gophermart/internal/logger"
	"github.com/timac11/yp-gophermart/internal/model"
)

type Worker struct {
	url                   string
	chOrdersForProcessing chan string
	chOrdersResult        chan model.Accrual
	client                *http.Client
	retryTimeout          time.Duration
	processingTimeout     time.Duration
}

func NewWorker(
	url string,
	chOrdersForProcessing chan string,
	chOrdersResult chan model.Accrual,
	timeoutConfig WorkerTimeoutsConfig,
) *Worker {
	return &Worker{
		url:                   url,
		chOrdersForProcessing: chOrdersForProcessing,
		chOrdersResult:        chOrdersResult,
		client:                &http.Client{Timeout: timeoutConfig.ClientTimeout},
		retryTimeout:          timeoutConfig.RetryTimeout,
		processingTimeout:     timeoutConfig.ProcessingTimeout,
	}
}

func (worker *Worker) Run(ctx context.Context) {
	log := logger.LoggerFromContext(ctx)

	for {
		select {
		case order := <-worker.chOrdersForProcessing:
			result, err := worker.processTask(ctx, order)
			if result != nil {
				worker.chOrdersResult <- *result
			} else {
				log.Error(err.Error())
				worker.chOrdersResult <- model.Accrual{Order: order, Status: model.AccrualInvalid}
			}
		case <-ctx.Done():
			return
		}
	}
}

func (worker *Worker) processTask(ctx context.Context, order string) (*model.Accrual, error) {
	var result *model.Accrual
	var err error

	retryCtx, cancel := context.WithTimeout(ctx, worker.processingTimeout)

	defer cancel()

	retry.Do(
		func() error {
			result, err = worker.executeRequest(order)
			return err
		},
		worker.getRetryOptions(retryCtx)...,
	)

	if result != nil {
		return result, nil
	}

	// by default accrual completed with error
	return nil, err
}

func (worker *Worker) executeRequest(order string) (*model.Accrual, error) {
	url := fmt.Sprintf("%s%s%s", worker.url, "/api/orders/", order)
	resp, err := worker.client.Get(url)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusOK {
		var result model.Accrual
		err = json.NewDecoder(resp.Body).Decode(&result)

		if err != nil {
			return nil, err
		}

		return &result, nil
	}

	if resp.StatusCode == http.StatusTooManyRequests {
		retryAfter, err := strconv.Atoi(resp.Header.Get("Retry-After"))

		if err != nil {
			return nil, err
		}

		return nil, errors.NewTooManyRequestsError(retryAfter)
	}

	return nil, errors.NewAccrualNotRegisteredError(order)
}

func (worker *Worker) getRetryOptions(ctx context.Context) []retry.Option {
	return []retry.Option{
		retry.Attempts(0),
		retry.RetryIf(func(err error) bool {
			var tooManyRequests *errors.TooManyRequestsError
			if _errors.As(err, &tooManyRequests) {
				return true
			}

			return false
		}),
		retry.DelayType(func(n uint, err error, config *retry.Config) time.Duration {
			var tooManyRequests *errors.TooManyRequestsError
			if _errors.As(err, &tooManyRequests) {
				return time.Duration(tooManyRequests.RetryAfter) * time.Second
			}

			return worker.retryTimeout
		}),
		retry.Context(ctx),
	}
}
