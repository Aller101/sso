package apps

import (
	"log/slog"
	"time"

	"ginhub.com/Aller101/sso/internal/apps/grpcapp"
)

type Apps struct {
	GRPCSrv *grpcapp.GRPCApp
}

func New(log *slog.Logger, grpcPort int, storagePath string, tokenTTL time.Duration) *Apps {

	grpcApp := grpcapp.New(log, grpcPort)

	return &Apps{
		GRPCSrv: grpcApp,
	}
}
