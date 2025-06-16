package apps

import (
	"log/slog"
	"time"

	"ginhub.com/Aller101/sso/internal/apps/grpcapp"
	"ginhub.com/Aller101/sso/internal/services/auth"
	"ginhub.com/Aller101/sso/internal/storage/sqlite"
)

type Apps struct {
	GRPCSrv *grpcapp.GRPCApp
}

func New(log *slog.Logger, grpcPort int, storagePath string, tokenTTL time.Duration) *Apps {

	storage, err := sqlite.New(storagePath)
	if err != nil {
		panic(err)
	}

	authService := auth.New(log, storage, tokenTTL)

	grpcApp := grpcapp.New(log, grpcPort)

	return &Apps{
		GRPCSrv: grpcApp,
	}
}
