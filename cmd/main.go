package main

import (
	"auth_service/internal/bootstrap"
	"auth_service/internal/config"
	"context"
	"github.com/joho/godotenv"
	"github.com/sethvargo/go-envconfig"
	"log"
	"os"
	"os/signal"
)

func main() {
	// 1. Загружаем .env в окружение
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  .env file not found, using system environment")
	}

	var cfg config.Config

	// 2. Читаем переменные окружения
	err := envconfig.ProcessWith(context.TODO(), &envconfig.Config{
		Target:   &cfg,
		Lookuper: envconfig.OsLookuper(),
	})
	if err != nil {
		panic(err)
	}

	// 3. Остальной код без изменений
	quitSignal := make(chan os.Signal, 1)
	signal.Notify(quitSignal, os.Interrupt)

	app := bootstrap.New(cfg)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-quitSignal
		cancel()
	}()

	app.Run(ctx)
}
