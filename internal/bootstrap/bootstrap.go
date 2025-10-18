package bootstrap

import (
	"auth_service/internal/adapter/driven/dbstore"
	"auth_service/internal/config"
	"auth_service/internal/usecase"
	"context"
	"fmt"
)

func initLayers(cfg config.Config) *App {
	teardown := make([]func(), 0)
	//log := logger.New(cfg.LogLevel, config.ServiceLabel, zap.WithCaller(true))

	db, err := initDB(*cfg.Postgres, "hpay_astrasend")
	if err != nil {
		panic(err)
	}

	storage := dbstore.New(db)
	//log.Info("Database connection established")

	teardown = append(teardown, func() {
		if err := db.Close(); err != nil {
			fmt.Println(err)
			//log.Error(err.Error())
		}
	})

	uc := usecase.New(cfg, storage)

	httpSrv := initHTTPService(&cfg, uc)

	teardown = append(teardown,
		func() {
			//log.Info("HTTP is shutting down")
			ctxShutDown, cancel := context.WithTimeout(context.Background(), gracefulDeadline)
			defer cancel()
			if err := httpSrv.Shutdown(ctxShutDown); err != nil {
				//log.Error(fmt.Sprintf("server Shutdown Failed:%s", err))
				return
			}
			//log.Info("HTTP is shut down")
		},
	)

	return &App{
		cfg:      cfg,
		rest:     httpSrv,
		teardown: teardown,
	}
}
