package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"ginhub.com/Aller101/sso/internal/apps"
	"ginhub.com/Aller101/sso/internal/config"
)

const (
	envLocal  = "local"
	envProd   = "prod"
	envDeploy = "deploy"
)

func main() {

	//go run .\cmd\sso\main.go --config=./config/local.yaml

	cfg := config.MustLoad()
	log := setupLogger(cfg.Env)

	log.Info("start sso", slog.Any("cfg", cfg))

	//TODO мб добавить структуру, где будет поле wg и можно будет увел или умен счетчик у wg в методах
	// wg := sync.WaitGroup{}

	applications := apps.New(log, cfg.Port, cfg.StoragePath, cfg.TokenTTL)

	// wg.Add(1)
	go applications.GRPCSrv.MustRun()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	<-sigCh
	applications.GRPCSrv.Stop()

	// wg.Wait()

}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	case envDeploy:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
