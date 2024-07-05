package repository

import "go.uber.org/fx"

func New() fx.Option {
	return fx.Module("repository",
		fx.Provide(
			postgres.New,
		),
		fx.Invoke(
			func(lc fx.Lifecycle, a *postgres.Repository) {
				lc.Append(fx.Hook{
					OnStart: a.OnStart,
					OnStop:  a.OnStop,
				})
			},
		),
	)
}
