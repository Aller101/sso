package app

import (
	"log/slog"
	"time"

	"ginhub.com/Aller101/sso/internal/app/grpcapp"
)

type App struct {
	GRPCSrv *grpcapp.GRPCApp
}

func New(log *slog.Logger, grpcPort int, storagePath string, tokenTTL time.Duration) *App {

	grpcApp := grpcapp.New(log, grpcPort)

	return &App{
		GRPCSrv: grpcApp,
	}
}
