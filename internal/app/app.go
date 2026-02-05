package app

import (
	"context"
	"net/http"

	accrualAgent "github.com/timac11/yp-gophermart/internal/accrual-agent"
	"github.com/timac11/yp-gophermart/internal/config"
	"github.com/timac11/yp-gophermart/internal/handler"
	"github.com/timac11/yp-gophermart/internal/repository"
)

func RunApplication() error {
	conf := config.InitConfig()

	// init repository
	repo, err := repository.NewPgClient(conf.DatabaseURI)

	if err != nil {
		return err
	}

	// run agent
	agent := accrualAgent.NewAgent(
		repo,
		conf.AccrualSystemAddress,
		accrualAgent.WorkerTimeoutsConfig{
			ClientTimeout:     conf.ClientTimeout,
			RetryTimeout:      conf.RetryTimeout,
			ProcessingTimeout: conf.ProcessingTimeout,
		},
		accrualAgent.AgentLimits{
			WorkersCount: conf.WorkersCount,
			OrdersBuffer: conf.OrdersBuffer,
		},
	)

	go agent.Run(context.Background())

	// run server

	mux, err := handler.InitRouter(conf, repo)
	if err != nil {
		return err
	}

	err = http.ListenAndServe(conf.Address, mux)
	if err != nil {
		return err
	}

	return nil
}
