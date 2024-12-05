package middleware

import (
	"auth/config"
	"auth/internal/domain/auth/repository/postgres"
	"go.uber.org/zap"
)

type Middleware struct {
	cfg   *config.ConfigModel
	repo  *postgres.Repository
	log   *zap.Logger
	roles map[string]int
}

func NewMiddleware(cfg *config.ConfigModel, log *zap.Logger, repository *postgres.Repository) *Middleware {
	return &Middleware{
		cfg:  cfg,
		log:  log,
		repo: repository,
	}
}
