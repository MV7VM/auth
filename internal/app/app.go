package app

import (
	"auth/config"
	"auth/internal/domain/auth/delivery/grpc"
	"context"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

func New() *fx.App {
	return fx.New(
		fx.Options(
			//repository.New(),
			//usecase.New(),
			grpc.New(),
		),
		fx.Provide(
			context.Background,
			config.NewConfig,
		),
		fx.WithLogger(
			func(log *zap.Logger) fxevent.Logger {
				return &fxevent.ZapLogger{Logger: log}
			},
		),
	)
}
