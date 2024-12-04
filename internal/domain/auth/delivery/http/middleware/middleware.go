package middleware

import (
	"auth/config"
	"go.uber.org/zap"
)

type Middleware struct {
	cfg *config.ConfigModel
	log *zap.Logger
}

func NewMiddleware(cfg *config.ConfigModel, log *zap.Logger) *Middleware {
	return &Middleware{
		cfg: cfg,
		log: log,
	}
}
