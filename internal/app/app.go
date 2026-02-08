package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	accrualAgent "github.com/timac11/yp-gophermart/internal/accrual-agent"
	"github.com/timac11/yp-gophermart/internal/config"
	"github.com/timac11/yp-gophermart/internal/handler"
	"github.com/timac11/yp-gophermart/internal/logger"
	"github.com/timac11/yp-gophermart/internal/repository"

	"golang.org/x/sync/errgroup"
)

func RunApplication() {
	conf := config.InitConfig()

	g, appCtx := errgroup.WithContext(context.Background())
	defer appCtx.Err()

	log := logger.LoggerFromContext(appCtx)

	// init repository
	repo, err := repository.NewPgClient(conf.DatabaseURI)
	if err != nil {
		log.Fatal(err.Error())
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

	g.Go(func() error {
		agent.Run(appCtx)
		return nil
	})

	// run server
	mux, err := handler.InitRouter(conf, repo)
	if err != nil {
		log.Fatal(err.Error())
	}

	server := NewServer(conf.Address, mux)
	g.Go(func() error {
		return server.Start()
	})

	// graceful shutdown
	quitChan := make(chan os.Signal, 1)
	signal.Notify(quitChan, syscall.SIGINT, syscall.SIGTERM)
	g.Go(func() error {
		<-quitChan
		log.Info("graceful shutdown signal")

		if err = server.Stop(appCtx); err != nil {
			return err
		}

		return nil
	})

	if err := g.Wait(); err != nil {
		log.Fatal(err.Error())
	}
}
