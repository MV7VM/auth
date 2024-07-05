package usecase

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func New() fx.Option {
	return fx.Module(
		"usecase",
		fx.Provide(
			NewUsecase,
		),
		fx.Invoke(
			func(lc fx.Lifecycle, u *Usecase) {
				lc.Append(fx.Hook{
					OnStart: u.OnStart,
				})
			},
		),
		fx.Decorate(func(log *zap.Logger) *zap.Logger {
			return log.Named("usecase")
		}),
	)
}
