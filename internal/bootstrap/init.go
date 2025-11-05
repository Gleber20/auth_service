package bootstrap

import (
	"auth_service/internal/adapter/driven/amqp"
	http2 "auth_service/internal/adapter/driving/http"
	"auth_service/internal/config"
	"auth_service/internal/usecase"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

func initDB(cfg config.Postgres, name string) (*sqlx.DB, error) {
	connConfig, err := pgx.ParseConfig(cfg.ConnectionURL())
	if err != nil {
		return nil, err
	}

	connStr := stdlib.RegisterConnConfig(connConfig)

	db, err := sqlx.Connect("pgx", connStr)
	if err != nil {
		return db, err
	}

	// Connection configuration
	// more info here https://www.alexedwards.net/blog/configuring-sqldb
	db.SetMaxOpenConns(cfg.MaxOpenConnections)
	db.SetMaxIdleConns(cfg.MaxIdleConnections)
	db.SetConnMaxLifetime(cfg.ConnectionMaxLifetime)

	return db, nil
}

func initHTTPService(
	cfg *config.Config,
	uc *usecase.UseCases,
) *http.Server {
	return http2.New(
		cfg,
		uc,
	)
}

type Resources struct {
	AMQPProducer *amqp.Producer
}

func InitResources(cfg *config.Config) (*Resources, error) {
	producer, err := amqp.NewProducer(cfg.AMQP_URL, "user-registered")
	if err != nil {
		log.Println("⚠️ Failed to initialize AMQP producer:", err)
		return nil, err
	}

	return &Resources{
		AMQPProducer: producer,
	}, nil
}
